package taskpriorityqueue_test

import (
	"github.com/donutloop/taskqueue"
	"testing"
)

func TestTaskPriorityQueue(t *testing.T) {

	queue := taskpriorityqueue.New(2)

	t1 := &taskpriorityqueue.Task{
		Priority: 10,
		Do: func() error {
			return nil
		},
	}
	queue.Push(t1)

	t2 := &taskpriorityqueue.Task{
		Priority: 10,
		Do: func() error {
			return nil
		},
	}
	queue.Push(t2)

	for i := 0; i < queue.Len(); i++ {
		if err := queue.Execute(); err != nil {
			t.Error("Unexpected error")
		}
	}
}
