package main

import (
	"errors"

	gs "github.com/rdkr/gitsync/sync"

	"math/rand"
	"sync"
	"time"
)

func main() {

	statuses := 10
	delayTime := time.Millisecond //* 2000

	ui := gs.NewUI(true, false, false)

	wg := sync.WaitGroup{}
	wg.Add(statuses + 2)

	go func() {
		ui.Run()
		wg.Done()
	}()
	go func() {
		for i := 0; i < statuses; i++ {
			time.Sleep(delayTime)
			ui.StatusChan <- getStatus(i)
			wg.Done()
		}
		close(ui.StatusChan)
		wg.Done()
	}()

	wg.Wait()

}

func getStatus(i int) gs.Status {
	statuses := []gs.Status{
		{"testpath1", gs.StatusError, "", errors.New("o no")},
		{"testpath2", gs.StatusFetched, "", nil},
		{"testpath3", gs.StatusCloned, "", nil},
		{"testpath4", gs.StatusUpToDate, "", nil},
	}
	s := rand.NewSource(int64(1337 + i))
	r := rand.New(s)
	return statuses[r.Intn(len(statuses))]
}
