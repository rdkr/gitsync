package cmd_test

import (
	"gitsync/cmd"
	"gitsync/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	git "gopkg.in/src-d/go-git.v4"
	git_memory "gopkg.in/src-d/go-git.v4/storage/memory"
)

func TestCloner(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCloner := mocks.NewMockCloner(ctrl)

	mockCloner.EXPECT().PlainOpen().Return(git.Init(git_memory.NewStorage(), nil))

	cmd.Clone(mockCloner, "somewhere")

}
