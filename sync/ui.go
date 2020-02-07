package sync

import (
	"fmt"
	"github.com/rdkr/gitsync/concurrency"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/gosuri/uilive"
)

const (
	SymbolError    = "\u001b[31m✘ \u001b[0m "
	SymbolClone    = "\u001b[36m✚  \u001b[0m"
	SymbolFetch    = "\u001b[33m↓  \u001b[0m"
	SymbolUpToDate = "\u001b[32m✔ \u001b[0m "
)

type UI struct {
	verbose                                         bool
	writer                                          *uilive.Writer
	cloneCount, fetchCount, upToDateCount, errCount int
	StatusChan                                      chan concurrency.Status
	statuses                                        []concurrency.Status
}

func ShouldBeVerbose(isTerminal, verbose, debug bool) bool {
	return !isTerminal || verbose || debug
}

func NewUI(isTerminal, verbose, debug bool) UI {

	verbose = ShouldBeVerbose(isTerminal, verbose, debug)

	writer := uilive.New() // TODO this is created even though its not necessarily used
	if !verbose {
		writer.Start()
	}

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return UI{
		verbose:       verbose,
		writer:        writer,
		cloneCount:    0,
		fetchCount:    0,
		upToDateCount: 0,
		errCount:      0,
		StatusChan:    make(chan concurrency.Status),
		statuses:      []concurrency.Status{},
	}
}

func (ui *UI) MakeUI(status concurrency.Status) string {
	var sb strings.Builder
	sb.WriteString("summary:")

	if status.Path != "" {
		ui.statuses = append(ui.statuses, status)
		if status.Err != nil {
			ui.errCount = ui.errCount + 1
		} else {
			switch status.Status {
			case concurrency.StatusCloned:
				ui.cloneCount = ui.cloneCount + 1
			case concurrency.StatusFetched:
				ui.fetchCount = ui.fetchCount + 1
			case concurrency.StatusUpToDate:
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

	return sb.String()
}

func (ui *UI) Run() {
	for {

		status, ok := <-ui.StatusChan
		if !ok {
			if !ui.verbose {
				ui.writer.Stop()
			}
			break
		}

		if !ui.verbose {
			switch status.Status {
			case concurrency.StatusCloned:
				_, err := fmt.Fprint(ui.writer, fmt.Sprintf(" %s%s\n", SymbolClone, status.Path))
				checkErr(err)
				ui.writer.Stop()
				ui.writer = uilive.New()
				ui.writer.Start()
			case concurrency.StatusError:
				_, err := fmt.Fprint(ui.writer, fmt.Sprintf(" %s%s - %s\n", SymbolError, status.Path, status.Err))
				checkErr(err)
				ui.writer.Stop()
				ui.writer = uilive.New()
				ui.writer.Start()
			}

			_, err := fmt.Fprint(ui.writer, ui.MakeUI(status))
			checkErr(err)
			err = ui.writer.Flush()
			checkErr(err)

		} else {
			fields := logrus.Fields{"path": status.Path}
			switch status.Status {
			case concurrency.StatusError:
				logrus.WithFields(fields).WithField("error", status.Err).Warn("error")
			case concurrency.StatusCloned:
				logrus.WithFields(fields).Info("cloned")
			case concurrency.StatusFetched:
				logrus.WithFields(fields).Debug("fetched")
			case concurrency.StatusUpToDate:
				logrus.WithFields(fields).Debug("up to date")
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
