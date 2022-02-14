package deleter

import (
	"context"
	"sync"
)

type Deleter struct {
	storage Storage
	upd     <-chan DeleteTask
	workers []chan DeleteTask
	queue   []DeleteTask
	mu      sync.Mutex
}

func New(storage Storage, upd <-chan DeleteTask) *Deleter {
	return &Deleter{
		storage: storage,
		upd:     upd,
		workers: []chan DeleteTask{},
		queue:   []DeleteTask{},
		mu:      sync.Mutex{},
	}
}

func (d *Deleter) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				for _, tasks := range d.workers {
					close(tasks)
				}
				return
			case task := <-d.upd:
				d.mu.Lock()
				d.queue = append(d.queue, task)
				d.mu.Unlock()
			}
		}
	}()
	var task DeleteTask
	for _, worker := range d.workers {
		if len(d.queue) > 0 {
			// queue not empty
			d.mu.Lock()
			task, d.queue = d.queue[0], d.queue[1:]
			worker <- task
			d.mu.Unlock()
		}
	}
}

func (d *Deleter) AddWorker() {
	workerChan := make(chan DeleteTask)
	d.workers = append(d.workers, workerChan)
	go func(workerChan chan DeleteTask) {
		for task := range workerChan {
			d.storage.BatchUpdateLinks(task)
		}
	}(workerChan)
}

type Storage interface {
	BatchUpdateLinks(task DeleteTask) error
}

type DeleteTask struct {
	UID   string
	Links []string
}
