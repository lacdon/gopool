package gopool

import (
	"errors"
	"sync"
)

// A FIFO data struct that only store Job, It CAN NOT be use without Init()
type Queue struct {
	sync.RWMutex
	sync.Cond

	jobSlice []Job
	stats    bool
}

// Queue CAN NOT be used without Init().
func (jq *Queue) Init() {
	jq.jobSlice = make([]Job, 0)

	// Queue is Composite by 2 parts: RWMutex and Cond.
	// Cond need a Locker to lock current area, so as Queue, so they need
	// SAME Locker which CAN NOT be passed by value, only can be passed by pointer.
	jq.L = jq

	jq.stats = true
}

// Put a new job into the end of the queue.
func (jq *Queue) Put(job Job) error {
	// Queue CAN NOT be used before Init() or after Destroy()
	if !jq.stats {
		return errors.New("put job into the queue, but queue stats is down")
	}
	if job == nil {
		return errors.New("put nil job into queue")
	}

	jq.Lock()
	defer jq.Unlock()
	jq.jobSlice = append(jq.jobSlice, job)

	// Notify block thread that we got a job now
	jq.Broadcast()
	return nil
}

// sync pop. If there is no job in the queue, it will block.
func (jq *Queue) Pop() (Job, error) {
	// Queue CAN NOT be used before Init() or after Destroy()
	if !jq.stats {
		return nil, errors.New("pop job out of the queue, but queue stats is down")
	}

	jq.Lock()
	defer jq.Unlock()

	// If there is no job in the queue, Pop() will block the thread and waiting for new job.
	for len(jq.jobSlice) == 0 && jq.stats {
		// Condition will doing those thing :
		//     1. Unlock mutex
		//     2. Waiting for signal, which emit by jq.Broadcast()
		//     3. Lock again
		jq.Wait()
	}

	// Check stats again in case of Queue has been Destroy()
	if !jq.stats {
		return nil, errors.New("pop job out of the queue, but queue stats is down")
	}

	// Get first element and move other elements forward
	ret := jq.jobSlice[0]
	jq.jobSlice = jq.jobSlice[1:]

	return ret, nil
}

// Return current Jobs number
func (jq *Queue) Size() int {
	jq.RLock()
	defer jq.RUnlock()

	return len(jq.jobSlice)
}

func (jq *Queue) Destroy() {
	jq.stats = false

	// Unblock all threads which are waiting for Cond's signal
	jq.Broadcast()
}
