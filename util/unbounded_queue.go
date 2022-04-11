package util

import (
	"sync"
)

type UnboundedQueue[T any] struct {
	outbound     []T
	outboundLock sync.RWMutex
	inbound      []T
	inboundLock  sync.RWMutex
}

func NewUnboundedQueue[T any]() *UnboundedQueue[T] {
	q := &UnboundedQueue[T]{}
	return q
}

func (q *UnboundedQueue[T]) EnQueue(items ...T) {
	q.inboundLock.Lock()
	defer q.inboundLock.Unlock()

	q.inbound = append(q.inbound, items...)

}

func (q *UnboundedQueue[T]) Consume(fn func([]T)) {
	q.outboundLock.Lock()
	defer q.outboundLock.Unlock()

	if len(q.outbound) == 0 {
		q.inboundLock.Lock()
		inboundLen := len(q.inbound)
		if inboundLen > 0 {
			q.inbound, q.outbound = q.outbound, q.inbound
		}
		q.inboundLock.Unlock()
	}

	if len(q.outbound) > 0 {
		fn(q.outbound)
		q.outbound = q.outbound[:0]
	}
}
