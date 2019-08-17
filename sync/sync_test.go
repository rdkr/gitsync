package sync_test

import (
	"errors"
	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"gitsync/mocks"
	"gitsync/sync"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"testing"
	"time"
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
			mockGit.EXPECT().PlainClone().Return("", nil)
			return mockGit
		},
		expected: sync.Status{"somewhere", "cloned", "", nil},
	},
	{
		name: "cloneFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {
			mockGit.EXPECT().PlainOpen().Return(nil, git.ErrRepositoryNotExists)
			mockGit.EXPECT().PlainClone().Return("", errors.New("uh oh"))
			return mockGit
		},
		expected: sync.Status{"somewhere", "", "", errors.New("unable to clone repo: uh oh")},
	},
	{
		name: "openFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {
			mockGit.EXPECT().PlainOpen().Return(nil, errors.New("uh oh"))
			return mockGit
		},
		expected: sync.Status{"somewhere", "", "", errors.New("unable to open repo: uh oh")},
	},
	{
		name: "headFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			// a bare repo with no head
			r, _ := git.Init(memory.NewStorage(), nil)

			mockGit.EXPECT().PlainOpen().Return(r, nil)

			return mockGit

		},
		expected: sync.Status{"somewhere", "", "", errors.New("unable to get head: reference not found")},
	},
	{
		name: "fetchSuccess",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			r := newRepo()

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			err = w.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName("test"), Create: true})
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().Fetch(r).Return("", git.NoErrAlreadyUpToDate)

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: "", Err: errors.New("not on master branch but fetched")},
	},
	{
		name: "fetchFail",
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {

			r := newRepo()

			w, err := r.Worktree()
			if err != nil {
				panic(err)
			}

			err = w.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName("test"), Create: true})
			if err != nil {
				panic(err)
			}

			mockGit.EXPECT().PlainOpen().Return(r, nil)
			mockGit.EXPECT().Fetch(r).Return("", errors.New("uh oh"))

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: "", Err: errors.New("not on master branch and: uh oh")},
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
			mockGit.EXPECT().Pull(w).Return("", nil)

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: "fetched", Err: nil},
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
			mockGit.EXPECT().Pull(w).Return("", git.NoErrAlreadyUpToDate)

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: "uptodate", Err: nil},
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
			mockGit.EXPECT().Pull(w).Return("", errors.New("uh oh"))

			return mockGit

		},
		expected: sync.Status{Path: "somewhere", Status: "", Err: errors.New("unable to pull master: uh oh")},
	},
}

func TestSync(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range syncTests {
		t.Run(tc.name, func(t *testing.T) {

			mockGit := mocks.NewMockGit(ctrl)
			mockGit = tc.setup(mockGit)

			actual := sync.GitSync(mockGit, "somewhere")

			if diff := deep.Equal(actual, tc.expected); diff != nil {
				t.Error(diff)
			}

		})
	}
}

func newRepo() *git.Repository {

	memFS := memfs.New()
	memFS.Create("test")

	r, err := git.Init(memory.NewStorage(), memFS)
	if err != nil {
		panic(err) // t.Error(err)
	}

	w, err := r.Worktree()
	if err != nil {
		panic(err) // t.Error(err)
	}

	_, err = w.Commit("Test Add And Commit", &git.CommitOptions{Author: &object.Signature{
		Name:  "foo",
		Email: "foo@foo.foo",
		When:  time.Now(),
	}})
	if err != nil {
		panic(err) // t.Error(err)
	}

	return r
}
