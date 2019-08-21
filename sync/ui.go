package sync

import (
	"fmt"
	"strings"

	"github.com/gosuri/uilive"
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
			case "cloned":
				ui.cloneCount = ui.cloneCount + 1
			case "fetched":
				ui.fetchCount = ui.fetchCount + 1
			case "uptodate":
				ui.upToDateCount = ui.upToDateCount + 1
			}
		}
	}

	if ui.cloneCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d \u001b[36m+\u001b[0m  ", ui.cloneCount))
	}
	if ui.fetchCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d \u001b[33m⟳\u001b[0m  ", ui.fetchCount))
	}
	if ui.upToDateCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d \u001b[32m✔\u001b[0m  ", ui.upToDateCount))
	}
	if ui.errCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d \u001b[31m✘\u001b[0m  ", ui.errCount))
	}

	sb.WriteString("\n")

	for _, status := range ui.statuses {
		if status.Err != nil {
			sb.WriteString(fmt.Sprintf(" \u001b[31m✘\u001b[0m  %s - %s\n", status.Path, status.Err))
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
