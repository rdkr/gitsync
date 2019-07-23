package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gosuri/uilive"
	"golang.org/x/crypto/ssh/terminal"
)

type status struct {
	path      string
	operation string
	err       error
}

type ui struct {
	isTerminal                                      bool
	writer                                          *uilive.Writer
	cloneCount, fetchCount, upToDateCount, errCount int
	statusChan                                      chan status
	statuses                                        []status
}

func newUI() ui {

	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))

	writer := uilive.New() // TODO this is created even though its not necessarily used
	if isTerminal {
		writer.Start()
		fmt.Fprint(writer.Newline(), "getting root group... ")
	}

	return ui{
		isTerminal:    isTerminal,
		writer:        writer,
		cloneCount:    0,
		fetchCount:    0,
		upToDateCount: 0,
		errCount:      0,
		statusChan:    make(chan status),
		statuses:      []status{},
	}
}

func (ui *ui) makeUI(root string, status status) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("getting root group... %s\nprocessing projects... ", root))

	ui.statuses = append(ui.statuses, status)

	if status.err != nil {
		ui.errCount = ui.errCount + 1
	} else {
		switch status.operation {
		case "clone":
			ui.cloneCount = ui.cloneCount + 1
		case "fetch":
			ui.fetchCount = ui.fetchCount + 1
		case "uptodate":
			ui.upToDateCount = ui.upToDateCount + 1
		}
	}

	if ui.cloneCount > 0 {
		sb.WriteString(fmt.Sprintf("%d \u001b[32m⬇\u001b[0m", ui.cloneCount))
	}
	if ui.fetchCount > 0 {
		sb.WriteString(fmt.Sprintf("%d \u001b[32m↻\u001b[0m", ui.fetchCount))
	}
	if ui.upToDateCount > 0 {
		sb.WriteString(fmt.Sprintf("%d \u001b[32m✔\u001b[0m", ui.upToDateCount))
	}
	if ui.errCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d \u001b[31m✘\u001b[0m", ui.errCount))
	}

	sb.WriteString("\n")

	for _, status := range ui.statuses {
		if status.err != nil {
			sb.WriteString(fmt.Sprintf(" \u001b[31m✘\u001b[0m %s: %s\n", status.path, status.err))
		}
		// } else {
		// 	sb.WriteString(fmt.Sprintf(" \u001b[32m✔\u001b[0m %s\n", status.path))
		// }
	}

	return sb.String()
}

func (ui *ui) run() {
	for {

		status, ok := <-ui.statusChan
		if !ok {
			break
		}

		if ui.isTerminal {
			fmt.Fprint(ui.writer.Newline(), ui.makeUI("placeholder", status))
			ui.writer.Flush() // it randomly prints multiple lines without this
		}
	}
}
