package cmd

import (
	"errors"
	"fmt"
	"testing"
)

func TestUI(t *testing.T) {
	ui := newUI()
	ui.makeUI(status{"test", "error", "", errors.New("err")})
	ui.makeUI(status{"test", "fetched", "", nil})
	ui.makeUI(status{"test", "cloned", "", nil})
	fmt.Println(ui.makeUI(status{"", "uptodate", "", nil}))
}
