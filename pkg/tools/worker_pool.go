package tools

import "sync"

type WorkerPool[T any] struct {
	work  chan *T
	wg    sync.WaitGroup
	count int
}

func NewWorkerPool[T any](count int, workFn func(T)) *WorkerPool[T] {
	wp := new(WorkerPool[T])
	wp.work = make(chan *T)
	wp.wg.Add(count)
	wp.count = count
	// Start the workers
	for i := 0; i < count; i++ {
		go func() {
			defer wp.wg.Done()
			for {
				next := <-wp.work
				if next == nil {
					return
				}
				workFn(*next)
			}
		}()
	}
	return wp
}

func (wp *WorkerPool[T]) Submit(item T) {
	wp.work <- &item
}

func (wp *WorkerPool[T]) Close() {
	for i := 0; i < wp.count; i++ {
		wp.work <- nil
	}
	wp.wg.Wait()
}
