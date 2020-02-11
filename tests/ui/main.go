package main

import (
	"errors"
	"github.com/rdkr/gitsync/concurrency"
	gs "github.com/rdkr/gitsync/sync"

	"math/rand"
	"sync"
	"time"
)

func main() {

	statuses := 200
	ui := gs.NewUI(true, false, false)

	wg := sync.WaitGroup{}
	wg.Add(statuses + 2)

	go func() {
		ui.Run()
		wg.Done()
	}()
	go func() {
		for i := 0; i < statuses; i++ {
			time.Sleep(time.Millisecond * 10)
			ui.StatusChan <- getStatus(i)
			wg.Done()
		}
		close(ui.StatusChan)
		wg.Done()
	}()

	wg.Wait()

}

func getStatus(i int) concurrency.Status {
	statuses := []concurrency.Status{
		{"testpath1", concurrency.StatusError, "", errors.New("o no")},
		{"testpath2", concurrency.StatusFetched, "", nil},
		{"testpath3", concurrency.StatusCloned, "", nil},
		{"testpath4", concurrency.StatusUpToDate, "", nil},
	}
	s := rand.NewSource(int64(1337 + i))
	r := rand.New(s)
	return statuses[r.Intn(len(statuses))]
}
