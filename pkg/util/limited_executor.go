package util

import "sync"

type LimitedConcurrentExecutor struct {
	limit chan int
	wg    *sync.WaitGroup
}

func NewLimitedConcurrentExecutor(limit int) *LimitedConcurrentExecutor {
	return &LimitedConcurrentExecutor{
		limit: make(chan int, limit),
		wg:    &sync.WaitGroup{},
	}
}

func (lce *LimitedConcurrentExecutor) add() {
	lce.limit <- 0
	lce.wg.Add(1)
}

func (lce *LimitedConcurrentExecutor) done() {
	<-lce.limit
	lce.wg.Done()
}

func (lce *LimitedConcurrentExecutor) Wait() {
	lce.wg.Wait()
}

func (lce *LimitedConcurrentExecutor) Execute(job func()) {
	lce.add()
	go func() {
		defer lce.done()
		job()
	}()
}
