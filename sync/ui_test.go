package sync_test

import (
	"errors"
	gs "github.com/rdkr/gitsync/sync"
	"sync"
	"testing"
)

func TestShouldBeVerbose(t *testing.T) {
	type args struct {
		isTerminal bool
		verbose    bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "script with verbose set",
			args: args{false, true},
			want: true,
		},
		{
			name: "script with verbose not set",
			args: args{false, false},
			want: true,
		},
		{
			name: "terminal with verbose set",
			args: args{true, true},
			want: true,
		},
		{
			name: "terminal without verbose not set",
			args: args{true, false},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gs.ShouldBeVerbose(tt.args.isTerminal, tt.args.verbose, false); got != tt.want {
				t.Errorf("ShouldBeVerbose() = %v, want %v", got, tt.want)
			}
		})
	}
}

// not really tests, but useful for checking the UI output formatters

func TestPrettyUI(t *testing.T) {
	testUI(false, false)
}

func TestVerboseUI(t *testing.T) {
	testUI(true, false)
}

func TestVerboseUIDebug(t *testing.T) {
	testUI(true, true)
}

func testUI(verbose, debug bool) {

	ui := gs.NewUI(true, verbose, debug)
	statuses := getStatuses()

	wg := sync.WaitGroup{}
	wg.Add(len(statuses) + 2)

	go func() {
		ui.Run()
		wg.Done()
	}()
	go func() {
		for _, status := range statuses {
			ui.StatusChan <- status
			wg.Done()
		}
		close(ui.StatusChan)
		wg.Done()
	}()

	wg.Wait()

}

func getStatuses() []gs.Status {
	return []gs.Status{
		{"testpath1", gs.StatusError, "", errors.New("o no")},
		{"testpath2", gs.StatusFetched, "", nil},
		{"testpath3", gs.StatusCloned, "", nil},
		{"testpath4", gs.StatusUpToDate, "", nil},
	}
}
