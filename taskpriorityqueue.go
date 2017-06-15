package taskpriorityqueue

import "reflect"

type Func interface{}

type Task struct {
	Arguments []interface{}
	Do        Func
	Priority  int64
	Index     int
}

type TaskPriorityQueue struct {
	queue []*Task
}

func New(capacity int) *TaskPriorityQueue {
	return &TaskPriorityQueue{
		queue: make([]*Task, 0, capacity),
	}
}

func (pq TaskPriorityQueue) Len() int {
	return len(pq.queue)
}

func (pq TaskPriorityQueue) Less(i, j int) bool {
	return pq.queue[i].Priority < pq.queue[j].Priority
}

func (pq TaskPriorityQueue) Swap(i, j int) {
	pq.queue[i], pq.queue[j] = pq.queue[j], pq.queue[i]
	pq.queue[i].Index = i
	pq.queue[j].Index = j
}

func (pq *TaskPriorityQueue) Push(t *Task) {
	validateFunc(t.Do, t.Arguments)

	n := len(pq.queue)
	c := cap(pq.queue)
	if n+1 > c {
		npq := make([]*Task, n, c*2)
		copy(npq, pq.queue)
		pq.queue = npq
	}
	pq.queue = append(pq.queue, t)
}

func (pq *TaskPriorityQueue) Execute() error {
	t := pq.pop()

	var outputVals []reflect.Value
	if len(t.Arguments) > 0 {
		arguments := []reflect.Value{reflect.ValueOf(t.Arguments)}
		outputVals = reflect.ValueOf(t.Do).Call(arguments)
	} else {
		outputVals = reflect.ValueOf(t.Do).Call(nil)
	}

	if err, ok := outputVals[0].Interface().(error); ok && err != nil {
		return err
	}
	return nil
}

func (pq *TaskPriorityQueue) pop() *Task {
	n := len(pq.queue)
	c := cap(pq.queue)
	if n < (c/2) && c > 25 {
		npq := make([]*Task, n, c/2)
		copy(npq, pq.queue)
		pq.queue = npq
	}
	t := pq.queue[n-1]
	t.Index = -1
	pq.queue = pq.queue[0 : n-1]
	return t
}

func validateFunc(fn Func, arguments []interface{}) {
	fnType := reflect.TypeOf(fn)

	// panic if conditions not met (because it's a programming error to have that happen)
	switch {
	case fnType.Kind() != reflect.Func:
		panic("value must be a function")
	case fnType.NumOut() != 1:
		panic("value must take exactly one output argument")
	case fnType.NumIn() != len(arguments):
		panic("Invalid of number of input arguments")
	}

	for i := 0; i < fnType.NumIn(); i++ {
		if fnType.In(i).Kind() != reflect.TypeOf(arguments[i]).Kind() {
			panic("argument doesn't matched")
		}
	}

	outType := fnType.Out(0)
	if ok := outType.Implements(reflect.TypeOf((*error)(nil)).Elem()); !ok {
		panic("func output argument must be a error")
	}
}
