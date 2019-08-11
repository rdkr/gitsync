package cmd_test

import (
	"fmt"
	"gitsync/cmd"
	"gitsync/mocks"
	//"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4/config"

	//"gopkg.in/src-d/go-git.v4/plumbing"
	//"gopkg.in/src-d/go-git.v4/plumbing/format/config"

	"testing"
	"github.com/golang/mock/gomock"
	"gopkg.in/src-d/go-git.v4"
	//"gopkg.in/src-d/go-git.v4/config"
	git_memory "gopkg.in/src-d/go-git.v4/storage/memory"
)


var syncTests = []struct {
	//in  string
	//out string
	setup func(*mocks.MockGit) *mocks.MockGit
}{
	{
		setup: func(mockGit *mocks.MockGit) *mocks.MockGit {


			r, _ := git.Init(git_memory.NewStorage(), nil)
			testBranch := &config.Branch{
				Name:   "foo",
				Remote: "origin",
				Merge:  "refs/heads/foo",
			}
			r.CreateBranch(testBranch)
			fmt.Println(r.Head())


			//repo, _ := git.Init(git_memory.NewStorage(), memfs.New())
			//
			//fmt.Println(repo.Worktree())


			//head, _ := repo.Head()
			//
			//
			//var ref plumbing.ReferenceName
			//ref = (plumbing.ReferenceName)("refs/heads/master")
			//
			//wt.Checkout(&git.CheckoutOptions{
			//	Branch: ref,
			//	Hash:   head.Hash(),
			//	Create: true,
			//})





			//branchName := fmt.Sprintf("refs/heads/%s", config.GitBranch)
			//branch := plumbing.ReferenceName(branchName)
			//worktree.Checkout(&git.CheckoutOptions{Create: true, Force: false, Branch: branch})


			//fmt.Println(repo.Head())

			//// First try to checkout branch
			//err = worktree.Checkout(&git.CheckoutOptions{Create: false, Force: false, Branch: b} )
			//
			//if err != nil {
			//	// got an error  - try to create it
			//	err :=  )
			//	CheckIfError(err)
			//}


			//repo.CreateBranch(git_config.Branch{
			//})

			//mockGit.EXPECT().PlainOpen().Return(repo, nil)

			return mockGit

		},
	},
}


func TestSync(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range syncTests {
		t.Run("tst", func(t *testing.T) {

			mockGit := mocks.NewMockGit(ctrl)
			mockGit = tc.setup(mockGit)

			expected := cmd.Status{Path: "somewhere", Status: "cloned", Err: nil}
			actual := cmd.Sync(mockGit, "somewhere")

			if expected != actual {
				t.Errorf("expected: %+v, got %+v", expected, actual)
			}

			//if s != tt.out {
			//	t.Errorf("got %q, want %q", s, tt.out)
			//}


		})
	}

}

