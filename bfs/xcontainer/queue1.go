package xcontainer

import (
	"container/list"
	"sync"
)

// Queue is a queue
type Queue interface {
	Front() *list.Element
	Len() int
	Add(interface{})
	Remove()
}

type queueImpl struct {
	*list.List
	mu sync.Mutex
}

func (q *queueImpl) Add(v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.PushBack(v)
}

func (q *queueImpl) Remove() {
	q.mu.Lock()
	defer q.mu.Unlock()
	e := q.Front()
	q.List.Remove(e)
}

// New is a new instance of a Queue
func New() Queue {
	return &queueImpl{list.New(), sync.Mutex{}}
}
