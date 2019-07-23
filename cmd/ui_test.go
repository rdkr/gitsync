package cmd

import (
	"errors"
	"fmt"
	"testing"
)

func TestUI(t *testing.T) {
	ui := newUI()
	ui.makeUI("test", status{"test", "error", "", errors.New("err")})
	ui.makeUI("test", status{"test", "fetched", "", nil})
	ui.makeUI("test", status{"test", "cloned", "", nil})
	fmt.Println(ui.makeUI("test", status{"test", "uptodate", "", nil}))
}
