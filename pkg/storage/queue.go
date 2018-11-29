package storage

import (
	"container/list"
)

type Queue struct {
	lists *list.List
}

func NewQueue() *Queue {
	return &Queue{
		lists: list.New(),
	}
}

func (q *Queue) Push(data []byte) {
	q.lists.PushBack(data)
}

func (q *Queue) Pop() []byte {
	element := q.lists.Front()
	if element == nil {
		return nil
	}

	q.lists.Remove(element)
	return element.Value.([]byte)
}
