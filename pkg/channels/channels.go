package channels

import (
	"os"
)

type Result struct {
	WorkerNumber int
	Value        int
}

type ChanGroups struct {
	Results chan Result
	Errors  chan error
	Done    chan os.Signal
}

func NewChanGroups(done chan os.Signal) *ChanGroups {
	return &ChanGroups{
		Results: make(chan Result),
		Errors:  make(chan error),
		Done:    done,
	}
}
