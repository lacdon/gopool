package gopool

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestWorker_Init(t *testing.T) {
	w := new(Worker)
	w.Init()
}

func TestWorker_Do(t *testing.T) {
	w := new(Worker)
	w.Init()

	wg := sync.WaitGroup{}
	loopTimes := 100
	wg.Add(loopTimes)
	for i := 0; i < loopTimes; i++ {
		temp := i
		f := func() {
			fmt.Printf("Test :job id=%d, time=%s\n", temp, time.Now())
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
			wg.Done()
		}
		e := w.Do(f)
		if e != nil {
			t.Error(e)
		}
	}
	wg.Wait()
}

// This test should be panic
func TestWorker_Kill(t *testing.T) {
	w := new(Worker)
	w.Init()

	wg := sync.WaitGroup{}
	loopTimes := 100
	wg.Add(loopTimes)
	for i := 0; i < loopTimes; i++ {
		temp := i
		f := func() {
			time.Sleep(time.Second)
			fmt.Printf("Test :job id=%d, time=%s\n", temp, time.Now())
			wg.Done()
		}
		e := w.Do(f)
		if e != nil {
			t.Error(e)
		}
	}
	go func() {
		time.Sleep(time.Second * 10)
		w.Kill()
	}()
	wg.Wait()
}
