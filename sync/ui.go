package sync

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gosuri/uilive"
)

const (
	SymbolError     = "\u001b[31m✗  \u001b[0m" //red
	SymbolClone     = "\u001b[36m+  \u001b[0m" //cyan
	SymbolFetch     = "\u001b[33m↓  \u001b[0m" //yellow
	SymbolUpToDate  = "\u001b[32m✓  \u001b[0m" //green
	SymbolUnmanaged = "\u001b[33m!  \u001b[0m"
)

type UI struct {
	verbose                                                         bool
	writer                                                          *uilive.Writer
	cloneCount, fetchCount, upToDateCount, errCount, unmanagedCount int
	StatusChan                                                      chan Status
	statuses                                                        []Status
	startTime                                                       int64
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
		verbose:        verbose,
		writer:         writer,
		cloneCount:     0,
		fetchCount:     0,
		upToDateCount:  0,
		errCount:       0,
		unmanagedCount: 0,
		StatusChan:     make(chan Status),
		statuses:       []Status{},
		startTime:      time.Now().UTC().UnixNano(),
	}
}

func (ui *UI) UpdateUI(status Status) {
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
			case StatusUnmanaged:
				ui.unmanagedCount = ui.unmanagedCount + 1
			}
		}
	}
}

func (ui *UI) MakeUI(done bool) string {
	var sb strings.Builder

	elapsed := time.Now().UTC().UnixNano() - ui.startTime
	timer := ((elapsed / 10000000) / 4) % 4
	icon := []string{" ◐  ", " ◓  ", " ◑  ", " ◒  "}

	if !done {
		sb.WriteString(icon[timer])
	}

	sb.WriteString("summary:")

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
	if ui.unmanagedCount > 0 {
		sb.WriteString(fmt.Sprintf(" %d %s", ui.unmanagedCount, SymbolUnmanaged))
	}

	sb.WriteString("\n")

	return sb.String()
}

func (ui *UI) Run() {

	for {

		select {
		case status, ok := <-ui.StatusChan:
			if !ok {
				if !ui.verbose {
					_, err := fmt.Fprint(ui.writer, ui.MakeUI(true))
					checkErr(err)
					ui.writer.Stop()
				}
				return
			}
			if !ui.verbose {
				switch status.Status {
				case StatusCloned:
					_, err := fmt.Fprintf(ui.writer, " %s%s\n", SymbolClone, status.Path)
					checkErr(err)
					ui.writer.Stop()
					ui.writer = uilive.New()
					ui.writer.Start()
				case StatusError:
					_, err := fmt.Fprintf(ui.writer, " %s%s - %s\n", SymbolError, status.Path, status.Err)
					checkErr(err)
					ui.writer.Stop()
					ui.writer = uilive.New()
					ui.writer.Start()
				case StatusUnmanaged:
					_, err := fmt.Fprintf(ui.writer, " %s%s - %s\n", SymbolUnmanaged, status.Path, status.Err)
					checkErr(err)
					ui.writer.Stop()
					ui.writer = uilive.New()
					ui.writer.Start()
				}

				ui.UpdateUI(status)

			} else {
				fields := logrus.Fields{"path": status.Path}
				switch status.Status {
				case StatusError:
					logrus.WithFields(fields).Warn(status.Err)
				case StatusCloned:
					logrus.WithFields(fields).Info("cloned")
				case StatusFetched:
					logrus.WithFields(fields).Info("fetched")
				case StatusUpToDate:
					logrus.WithFields(fields).Info("up to date")
				case StatusUnmanaged:
					logrus.WithFields(fields).Warn("unmanaged")
				}
			}
		case <-time.After(10 * time.Millisecond):
		}
		if !ui.verbose {
			_, err := fmt.Fprint(ui.writer, ui.MakeUI(false))
			checkErr(err)
			err = ui.writer.Flush()
			checkErr(err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}
