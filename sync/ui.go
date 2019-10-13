package sync

import (
	"fmt"
	"strings"

	"github.com/gosuri/uilive"
)

const (
	SymbolError    = "\u001b[31m✘ \u001b[0m "
	SymbolClone    = "\u001b[36m🞦 \u001b[0m"
	SymbolFetch    = "\u001b[33m🡻 \u001b[0m"
	SymbolUpToDate = "\u001b[32m✔ \u001b[0m "
)

type UI struct {
	verbose                                         bool
	writer                                          *uilive.Writer
	cloneCount, fetchCount, upToDateCount, errCount int
	StatusChan                                      chan Status
	statuses                                        []Status
}

func ShouldBeVerbose(isTerminal, verbose bool) bool {
	return !isTerminal || verbose
}

func NewUI(isTerminal, verbose bool) UI {

	verbose = ShouldBeVerbose(isTerminal, verbose)

	writer := uilive.New() // TODO this is created even though its not necessarily used
	if !verbose {
		writer.Start()
	}

	return UI{
		verbose:       verbose,
		writer:        writer,
		cloneCount:    0,
		fetchCount:    0,
		upToDateCount: 0,
		errCount:      0,
		StatusChan:    make(chan Status),
		statuses:      []Status{},
	}
}

func (ui *UI) MakeUI(status Status) string {
	var sb strings.Builder
	sb.WriteString("result:")

	if status.Path != "" {
		ui.statuses = append(ui.statuses, status)
		if status.Err != nil {
			ui.errCount = ui.errCount + 1
		} else {
			switch status.Status {
			case StatusCloned:
				ui.cloneCount = ui.cloneCount + 1
			case StatusFetched:
				ui.fetchCount = ui.fetchCount + 1
			case StatusUpToDate:
				ui.upToDateCount = ui.upToDateCount + 1
			}
		}
	}

	if ui.cloneCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d %s", ui.cloneCount, SymbolClone))
	}
	if ui.fetchCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d %s", ui.fetchCount, SymbolFetch))
	}
	if ui.upToDateCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d %s", ui.upToDateCount, SymbolUpToDate))
	}
	if ui.errCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d %s", ui.errCount, SymbolError))
	}

	sb.WriteString("\n")

	for _, status := range ui.statuses {
		switch status.Status {
		case StatusCloned:
			sb.WriteString(fmt.Sprintf(" %s%s\n", SymbolClone, status.Path))
		case StatusError:
			sb.WriteString(fmt.Sprintf(" %s%s - %s\n", SymbolError, status.Path, status.Err))
		}
	}

	return sb.String()
}

func (ui *UI) Run() {
	for {

		status, ok := <-ui.StatusChan
		if !ok {
			break
		}

		if !ui.verbose {
			fmt.Fprint(ui.writer.Newline(), ui.MakeUI(status))
			ui.writer.Flush() // it randomly prints multiple lines without this
		} else {
			fmt.Println(status)
		}
	}
}
