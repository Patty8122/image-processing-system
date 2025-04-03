package concurrent

import (
	"sync"
	"sync/atomic"
	// "fmt"
	"errors"

)

type Task interface{}

// Bounded Deque implements lock-free FIFO freelist based queue.
type BDEQueue interface {
	PushBottom(task Task) error
	PopTop() Task
	PopBottom() Task
	Size() int
	IsEmpty() bool
}


type queue struct {
	top int64 // top stamp is the first 32 bits of top and top value is the last 32 bits of top
	bottom int64 // bottom stamp is the first 32 bits of bottom and bottom value is the last 32 bits of bottom
	tasks []Task
	size     int
	mtx      *sync.Mutex
	capacity int
}

func NewBoundedDEQueue(capacity int) BDEQueue {
	// capacity is the size of the queue and it cannot be more than 2^32
	// top and bottom are 32 bit integers
	// assert(capacity < 1<<32, "Capacity cannot be more than 2^32")
	queue := &queue{top: 0, bottom: 0, tasks: make([]Task, capacity), size: 0, mtx: &sync.Mutex{}, capacity: capacity}
	return queue
}

func (queue *queue) Size() int {
	queue.mtx.Lock()
	defer queue.mtx.Unlock()
	return queue.size
}

func (queue *queue) PopTop() Task {
	// lock free pop top
 	oldTop := atomic.LoadInt64(&queue.top)
	oldStamp := (oldTop >> 32)
	oldTop = (oldTop & 0xFFFFFFFF)
	newTop := oldTop + 1
	newStamp := oldStamp + 1

	if ( (queue.bottom & 0xFFFFFFFF) <= oldTop) {
		return nil
	}

	var t Task = queue.tasks[oldTop]
	if (atomic.CompareAndSwapInt64(&queue.top, (oldTop | (oldStamp << 32)), (newTop | (newStamp << 32)))) {
		queue.mtx.Lock()
		queue.size -= 1
		queue.mtx.Unlock()
		return t
	}
	return nil
}

func (queue *queue) PushBottom(task Task) error {
	// lock free push bottom
	// Steps:
	// Since only the bottom can push, we don't need to CAS the bottom
	// We can just push the task to the bottom and increment the bottom

	if (queue.bottom & 0xFFFFFFFF) > int64(queue.capacity) {
		// queue is full
		return errors.New("Full queue")

	}
	queue.tasks[(queue.bottom & 0xFFFFFFFF)] = task
	queue.bottom += 1 | 1 << 32 // increment bottom and stamp
	queue.mtx.Lock()
	queue.size += 1
	queue.mtx.Unlock()
	return nil
}

func (queue *queue) PopBottom() Task {
	// lock free pop bottom
	// Steps:
	// This would fail when the bottom is being popped by current thread and another thread is trying to steal
	// So we need to CAS the bottom
	// If CAS fails, the steal would succeed and we can simply return nil
	// If CAS succeeds, we can return the task at the bottom

	// check if queue is empty
	if ((queue.bottom & 0xFFFFFFFF) == 0) {
		return nil
	}

	queue.bottom -= 1 // decrement bottom since we are popping 
	var t Task = queue.tasks[queue.bottom & 0xFFFFFFFF]

	// adjust top
	oldTop := atomic.LoadInt64(&queue.top)
	oldStamp := (oldTop >> 32)
	oldTop = int64(oldTop & 0xFFFFFFFF)
	
	newTop := int64(0)
	newStamp := oldStamp + 1

	if ((queue.bottom & 0xFFFFFFFF) > oldTop) {
		// reduce size
		queue.mtx.Lock()
		queue.size -= 1
		queue.mtx.Unlock()
		return t
	}

	if ((queue.bottom & 0xFFFFFFFF) == oldTop) {
		queue.bottom = ((queue.bottom >> 32) << 32) // reset bottom to 0 while keeping the stamp
		if (atomic.CompareAndSwapInt64(&queue.top, (oldTop | (oldStamp << 32)), (newTop | (newStamp << 32)))) {
			queue.mtx.Lock()
			queue.size -= 1
			queue.mtx.Unlock()
			return t
		}
	}

	// if we reach here, we failed to pop bottom
	queue.top = ((newStamp << 32) | newTop)
	queue.bottom = ((queue.bottom >> 32) << 32) // reset bottom to 0 while keeping the stamp

	// decrement size
	queue.mtx.Lock()
	queue.size -= 1
	queue.mtx.Unlock()

	return nil
}

func (queue *queue) IsEmpty() bool {
	queue.mtx.Lock()
	defer queue.mtx.Unlock()
	if (queue.size == 0) {
		return true
	}
	return false

}
