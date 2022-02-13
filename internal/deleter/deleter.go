package deleter

import (
	"context"
)

type Deleter struct {
	storage Storage
	upd     <-chan DeleteTask
	workers []chan DeleteTask
}

func New(storage Storage, upd <-chan DeleteTask) *Deleter {
	return &Deleter{storage: storage, upd: upd}
}

func (d *Deleter) Run(ctx context.Context) {
	for _, worker := range d.workers {
		select {
		case <-ctx.Done():
			for _, tasks := range d.workers {
				close(tasks)
			}
			return
		case task := <-d.upd:
			worker <- task
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
