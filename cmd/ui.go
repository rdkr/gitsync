package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosuri/uilive"
	"golang.org/x/crypto/ssh/terminal"
)

type Status struct {
	Path   string
	Status string
	Output string
	Err    error
}

type ui struct {
	prettyPrint                                     bool
	writer                                          *uilive.Writer
	cloneCount, fetchCount, upToDateCount, errCount int
	statusChan                                      chan Status
	statuses                                        []Status
}

func NewUI(verbose bool) ui {

	prettyPrint := !verbose && terminal.IsTerminal(int(os.Stdout.Fd()))

	writer := uilive.New() // TODO this is created even though its not necessarily used
	if !verbose {
		writer.Start()
	}

	return ui{
		prettyPrint:   prettyPrint,
		writer:        writer,
		cloneCount:    0,
		fetchCount:    0,
		upToDateCount: 0,
		errCount:      0,
		statusChan:    make(chan Status),
		statuses:      []Status{},
	}
}

func (ui *ui) makeUI(status Status) string {
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

func (ui *ui) run() {
	for {

		status, ok := <-ui.statusChan
		if !ok {
			break
		}

		if ui.prettyPrint {
			fmt.Fprint(ui.writer.Newline(), ui.makeUI(status))
			ui.writer.Flush() // it randomly prints multiple lines without this
		} else {
			fmt.Println(status)
		}
	}
}
