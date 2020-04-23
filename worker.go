package gopool

import (
	"errors"
	"fmt"
)

// One Worker need 2 go routines :
//   - For do current job
//   - For receive new job from jobC and put job into jobQ
type Worker struct {
	jobQ  Queue    // For store all job
	jobC  chan Job // For receive new job
	stats bool     // stats flag
}

func (w *Worker) Init() {
	w.jobQ.Init()
	w.jobC = make(chan Job)

	w.stats = true

	go w.doJob()
	go w.recvJob()
}

func (w *Worker) Kill() {
	w.stats = false
	close(w.jobC)

	w.jobQ.Destroy()
}

func (w Worker) Stats() bool {
	return w.stats
}

// Worker will not do new job immediately.
// Worker will do the oldest job in the order of they are received.
// This function will not block thread
func (w *Worker) Do(job Job) error {
	if job == nil {
		return errors.New("job is nil")
	}

	if !w.stats {
		return errors.New("worker not ready")
	}

	w.jobC <- job
	return nil
}

// Waiting for a new job from queue and do it.
func (w *Worker) doJob() {
	for w.stats {
		// Pop the oldest job from the queue
		f, e := w.jobQ.Pop()
		if e != nil {
			fmt.Printf("Worker do Job error when pop job from queue. Err=%v\n", e)
			continue
		}

		// Do job here
		f()
	}
}

// get new job from channel and put into queue
func (w *Worker) recvJob() {
	for w.stats {
		job, ok := <-w.jobC
		// jobC is closed
		if !ok {
			w.stats = false
			return
		}

		e := w.jobQ.Put(job)
		if e != nil {
			fmt.Printf("Worker cecv Job error when put job into queue. Err=%v\n", e)
			continue
		}
	}
}
