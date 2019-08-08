package cmd

import (
	"errors"
	"fmt"
	"testing"
)

// not really a test, but useful for checking the UI output formater
func TestUI(t *testing.T) {
	ui := newUI(false)
	ui.makeUI(Status{"testpath1", "error", "", errors.New("err")})
	ui.makeUI(Status{"testpath2", "fetched", "", nil})
	ui.makeUI(Status{"testpath3", "cloned", "", nil})
	ui.makeUI(Status{"testpath4", "uptodate", "", nil})
	fmt.Println(ui.makeUI(Status{"done", "", "", nil}))
}
