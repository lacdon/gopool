package gopool

import (
	"fmt"
	"testing"
)

func TestQueue_Init(t *testing.T) {
	q := new(Queue)
	q.Init()
}

func TestQueue_Pop(t *testing.T) {
	q := new(Queue)
	q.Init()
	for i := 0; i < 100; i++ {
		temp := i
		q.Put(func() {
			fmt.Printf("%d\t", temp)
		})
	}
	for i := 0; i < 100; i++ {
		f, err := q.Pop()
		if err != nil {
			t.Error(err)
		}
		f()
	}

}

func TestQueue_Destroy(t *testing.T) {
	q := new(Queue)
	q.Init()
	for i := 0; i < 100; i++ {
		temp := i
		q.Put(func() {
			fmt.Printf("%d\t", temp)
		})
	}
	for i := 0; i < 100; i++ {
		f, err := q.Pop()
		if err != nil {
			t.Fatal(err)
		}
		f()
		if i == 50 {
			q.Destroy()
		}
	}
}
