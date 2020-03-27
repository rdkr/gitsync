package sync_test

import (
	"errors"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/rdkr/gitsync/concurrency"
	"github.com/rdkr/gitsync/mocks"
	"github.com/rdkr/gitsync/sync"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var syncTests = []struct {
	name     string
	setup    func(*mocks.MockGit) *mocks.MockGit
	expected sync.Status
}{
	{
		name: "cloneSuccess",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {
			mockGit.EXPECT().PlainOpen().Return(nil, git.ErrRepositoryNotExists)
			mockGit.EXPECT().GetLocation().Return("somewhere")
			mockGit.EXPECT().PlainClone().Return("", nil)
			return mockGit
		},
		expected: sync.Status{"somewhere", sync.StatusCloned, "", nil},
	},
	{
		name: "cloneFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {
			mockGit.EXPECT().PlainOpen().Return(nil, git.ErrRepositoryNotExists)
			mockGit.EXPECT().GetLocation().Return("somewhere")
			mockGit.EXPECT().PlainClone().Return("", errors.New("uh oh"))
			return mockGit
		},
		expected: sync.Status{"somewhere", sync.StatusError, "", errors.New("unable to clone repo: uh oh")},
	},
	{
		name: "openFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {
			mockGit.EXPECT().PlainOpen().Return(nil, errors.New("uh oh"))
			mockGit.EXPECT().GetLocation().Return("somewhere")
			return mockGit
		},
		expected: sync.Status{"somewhere", sync.StatusError, "", errors.New("unable to open repo: uh oh")},
	},
	{
		name: "workTreeFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			// a bare repo with no worktree
			r, err := git.Init(memory.NewStorage(), nil)
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().GetLocation().Return("somewhere")

			return mockGit

		},
		expected: sync.Status{"somewhere", sync.StatusError, "", errors.New("unable to get worktree: worktree not available in a bare repository")},
	},
	{
		name: "headFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			memFS := memfs.New()
			_, err := memFS.Create("test")
			if err != nil {
				panic(err)
			}

			r, err := git.Init(memory.NewStorage(), memFS)
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().GetLocation().Return("somewhere")

			return mockGit

		},
		expected: sync.Status{"somewhere", sync.StatusError, "", errors.New("unable to get head: reference not found")},
	},
	{
		name: "fetchSuccess",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			r := newRepo()

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			err = w.Checkout(&git.CheckoutOptions{Branch: "test", Create: true})
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().GetLocation().Return("somewhere")
			mockGit.EXPECT().Fetch(r).Return("", git.NoErrAlreadyUpToDate)

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: sync.StatusError, Err: errors.New("not on master branch but fetched")},
	},
	{
		name: "fetchFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			r := newRepo()

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			err = w.Checkout(&git.CheckoutOptions{Branch: "test", Create: true})
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().GetLocation().Return("somewhere")
			mockGit.EXPECT().Fetch(r).Return("", errors.New("uh oh"))

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: sync.StatusError, Err: errors.New("not on master branch and: uh oh")},
	},
	{
		name: "pullSuccess",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			r := newRepo()

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().GetLocation().Return("somewhere")
			mockGit.EXPECT().Pull(w).Return("", nil)

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: sync.StatusFetched, Err: nil},
	},
	{
		name: "pullUpToDate",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			r := newRepo()

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().GetLocation().Return("somewhere")
			mockGit.EXPECT().Pull(w).Return("", git.NoErrAlreadyUpToDate)

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: sync.StatusUpToDate, Err: nil},
	},
	{
		name: "pullFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			r := newRepo()

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().GetLocation().Return("somewhere")
			mockGit.EXPECT().Pull(w).Return("", errors.New("uh oh"))

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: sync.StatusError, Err: errors.New("unable to pull master: uh oh")},
	},
}

func TestSync(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range syncTests {
		t.Run(tc.name, func(t *testing.T) {

			mockGit := mocks.NewMockGit(ctrl)
			mockGit = tc.setup(mockGit)

			actual := sync.GitSync(mockGit)

			if diff := deep.Equal(actual, tc.expected); diff != nil {
				t.Error(diff)
			}

		})
	}
}

type mockGit struct {
	concurrency.Project
	*mocks.MockGit
}

func newRepo() *git.Repository {

	memFS := memfs.New()
	_, err := memFS.Create("test")
	if err != nil {
		panic(err)
	}

	r, err := git.Init(memory.NewStorage(), memFS)
	if err != nil {
		panic(err)
	}

	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	_, err = w.Commit("Test Add And Commit", &git.CommitOptions{Author: &object.Signature{
		Name:  "foo",
		Email: "foo@foo.foo",
		When:  time.Now(),
	}})
	if err != nil {
		panic(err)
	}

	return r
}
