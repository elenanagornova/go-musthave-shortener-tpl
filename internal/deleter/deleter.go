package deleter

import (
	"context"
	"sync"
)

type Deleter struct {
	storage Storage
	upd     <-chan DeleteTask
}

func New(storage Storage, upd <-chan DeleteTask) *Deleter {
	return &Deleter{
		storage: storage,
		upd:     upd,
	}
}

func (d *Deleter) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case task := <-d.upd:
			go func() {
				wg.Add(1)
				defer wg.Done()
				d.AddWorker(ctx, task)
			}()
		}
	}
}

func (d *Deleter) AddWorker(ctx context.Context, task DeleteTask) {
	d.storage.BatchUpdateLinks(ctx, task)
}

type Storage interface {
	BatchUpdateLinks(ctx context.Context, task DeleteTask) error
}

type DeleteTask struct {
	UID   string
	Links []string
}
