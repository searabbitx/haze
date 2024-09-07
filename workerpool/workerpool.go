package workerpool

import (
	"sync"
)

type Pool struct {
	wg    *sync.WaitGroup
	input chan func()
}

func worker(input chan func(), wg *sync.WaitGroup) {
	defer wg.Done()
	for f := range input {
		f()
	}
}

func NewPool(size int) Pool {
	wg := new(sync.WaitGroup)
	input := make(chan func())

	for i := 0; i < size; i++ {
		wg.Add(1)
		go worker(input, wg)
	}

	return Pool{wg, input}
}

func (p Pool) RunTask(task func()) {
	p.input <- task
}

func (p Pool) Wait() {
	close(p.input)
	p.wg.Wait()
}
