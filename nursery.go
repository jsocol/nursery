package nursery

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var ErrStopped = errors.New("nursery is stopping or stopped")

// A Nursery is a control flow mechanism for concurrent or parallel tasks.
type Nursery interface {
	Start(Task) error
}

// A Task is a job
type Task func(context.Context) error

type Initializer func(Nursery) error

// Open creates a new Nursery that can be used to start Tasks. All started
// Tasks are guaranteed to be completed by the time Open returns. In terms of
// the Go memory model, all Tasks "synchronize before" Open can complete.
func Open(ctx context.Context, init Initializer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var err error

	n := &nursery{
		ctx:     ctx,
		cancel:  cancel,
		cond:    sync.NewCond(&sync.Mutex{}),
		errChan: make(chan error, 1),
	}

	doneChan := make(chan struct{})
	go func() {
		n.cond.L.Lock()
		for !n.started.Load() || n.running.Load() > 0 {
			n.cond.Wait()
		}
		close(doneChan)
		n.stopping.Store(true)
		close(n.errChan)
	}()

	go func() {
		err = <-n.errChan
		if err != nil {
			n.stopping.Store(true)
			n.cancel()
		}
	}()

	n.started.Store(true)
	err = init(n)
	if err != nil {
		n.errChan <- err
	}

	<-doneChan
	if err != nil {
		return err
	}
	if n.ctx.Err() != nil {
		return n.ctx.Err()
	}
	return nil
}

type nursery struct {
	ctx      context.Context
	cancel   context.CancelFunc
	started  atomic.Bool
	running  atomic.Int64
	stopping atomic.Bool
	cond     *sync.Cond
	errChan  chan error
}

func (n *nursery) Start(task Task) error {
	if n.stopping.Load() {
		return ErrStopped
	}

	n.running.Add(1)

	go func() {
		defer n.running.Add(-1)
		defer n.cond.Signal()

		err := task(n.ctx)
		if err != nil && !n.stopping.Load() {
			n.errChan <- err
		}
	}()

	return nil
}
