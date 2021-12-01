package queue

import "errors"

var IsFullErr = errors.New("queue is full")

func NewStringQueue(size int) *StringQueue {
	return &StringQueue{
		c: make(chan string, size),
	}
}

type StringQueue struct {
	c chan string
}

func (q *StringQueue) Push(s string) error {
	select {
	case q.c <- s:
		return nil
	default:
		return IsFullErr
	}
}

func (q *StringQueue) Pop() (string, bool) {
	select {
	case s := <-q.c:
		return s, true
	default:
		return "", false
	}
}

func (q *StringQueue) Await() string {
	return <-q.c
}
