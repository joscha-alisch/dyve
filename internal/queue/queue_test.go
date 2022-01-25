package queue

import (
	"errors"
	"testing"
	"time"
)

func TestQueueSize(t *testing.T) {
	size := 2
	q := NewStringQueue(size)
	for i := 0; i < size; i++ {
		err := q.Push("a")
		if err != nil {
			t.Fatal("expected no error")
		}
	}
	err := q.Push("a")
	if !errors.Is(err, IsFullErr) {
		t.Fatal("expected IsFullError, but was ", err)
	}
}

func TestStringQueue_Pop(t *testing.T) {
	q := NewStringQueue(1)

	_ = q.Push("element")
	el, ok := q.Pop()
	if !ok {
		t.Fatal("expected element in queue")
	}

	if el != "element" {
		t.Fatal("popped item should be 'element' but was ", el)
	}
}

func TestStringQueue_PopEmpty(t *testing.T) {
	q := NewStringQueue(1)

	_, ok := q.Pop()
	if ok {
		t.Fatal("expected no element in queue")
	}
}

func TestStringQueue_Await(t *testing.T) {
	q := NewStringQueue(10)
	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = q.Push("element")
	}()

	res := q.Await()
	if res != "element" {
		t.Fatal("expected result to be 'element' but was ", res)
	}
}
