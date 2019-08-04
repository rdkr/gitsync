package cmd

import (
	"errors"
	"fmt"
	"testing"
)

func TestUI(t *testing.T) {
	ui := newUI(false)
	ui.makeUI(status{"testpath1", "error", "", errors.New("err")})
	ui.makeUI(status{"testpath2", "fetched", "", nil})
	ui.makeUI(status{"testpath3", "cloned", "", nil})
	ui.makeUI(status{"testpath4", "uptodate", "", nil})
	fmt.Println(ui.makeUI(status{"done", "", "", nil}))
}
