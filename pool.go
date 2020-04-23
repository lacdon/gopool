package gopool

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type Scheduler struct {
	sync.RWMutex

	jobChan     chan Job
	workersSize uint
	workers     map[uint]*Worker
	status      bool
}

func New(workerNum uint) (*Scheduler, error) {
	if workerNum <= 0 {
		return nil, errors.New("new Scheduler error, workerNum <= 0")
	}
	ret := new(Scheduler)
	ret.Init(workerNum)
	return ret, nil
}

func (s *Scheduler) Init(workerSize uint) {
	s.Lock()
	defer s.Unlock()

	s.workersSize = workerSize
	s.jobChan = make(chan Job)
	s.workers = make(map[uint]*Worker)

	for i := uint(0); i < s.workersSize; i++ {
		worker := new(Worker)
		worker.Init()
		s.workers[i] = worker
	}

	s.status = true
	go s.handleJob()
}

func (s *Scheduler) Dispatch(job Job) error {
	if !s.status {
		return errors.New("scheduler has been closed")
	}
	s.jobChan <- job
	return nil
}

// Close will delete all worker, but won't stop worker's current job
func (s *Scheduler) Close() {
	s.Lock()
	defer s.Unlock()

	s.status = false
	close(s.jobChan)
	for k, worker := range s.workers {
		if worker != nil {
			// If worker is doing a blocking job, it will continue do it.
			worker.Kill()
		}
		delete(s.workers, k)
	}
}

func (s *Scheduler) handleJob() {
	for s.status {
		job, ok := <-s.jobChan
		// chan has been closed cause by s.Close()
		if !ok {
			s.status = false
			return
		}

		if job == nil {
			continue
		}

		// Simple load balance by rand worker id
		workerId := uint(rand.Intn(int(s.workersSize) - 1))
		s.RLock()
		worker, ok := s.workers[workerId]
		s.RUnlock()

		// Worker should not be nil
		if !ok || worker == nil {
			s.status = false
			return
		}

		// Async do job by one worker
		err := worker.Do(job)
		if err != nil {
			fmt.Printf("worker id=%d can not do job err=%v", workerId, err)
		}
	}
}
