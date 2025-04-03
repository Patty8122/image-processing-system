package concurrent

import (
	"math/rand"
	"sync"
	"time"
)

type Exec interface {
	Submit(task interface{})
	Shutdown()
}

type WorkStealingExecutor struct {
	wg           *sync.WaitGroup
	capacity     int
	shutdown     *bool
	currIndex    int
	totalTasks   int
	localGoroutineQueues []BDEQueue
	mtx          *sync.Mutex
}

func NewWorkStealingExecutor(capacity, threshold int) Exec {
	taskQueues := make([]BDEQueue, capacity)
	for i := 0; i < capacity; i += 1 {
		queue := NewBoundedDEQueue(15000)
		taskQueues[i] = queue
	}

	executor := &WorkStealingExecutor{
		wg:           &sync.WaitGroup{},
		capacity:     capacity,
		shutdown:     new(bool),
		totalTasks:   0,
		localGoroutineQueues: taskQueues,
		mtx:          &sync.Mutex{},
	}
	executor.start()
	return executor
}

func (w *WorkStealingExecutor) start() {
	for i := 0; i < w.capacity; i += 1 {
		w.wg.Add(1)
		go stealingWorker(w, i)
	}
}

func stealingWorker(w *WorkStealingExecutor, threadIdx int) {
	defer w.wg.Done()
	// Loop until shutdown is true and all queues are empty
	for {
		if w.localGoroutineQueues[threadIdx].IsEmpty() {

				rand.Seed(time.Now().UnixNano())
				hostThread := threadIdx
				for hostThread == threadIdx {
					hostThread = rand.Intn(w.capacity)
				}

				if w.localGoroutineQueues[hostThread].Size() > 2 {
					// steal
					currTask, _ := (w.localGoroutineQueues[hostThread].PopTop()).(Request)
					processImage(currTask, currTask.dataDir)
					w.mtx.Lock()
					w.totalTasks -= 1
					w.mtx.Unlock()
				}
		
		} else {
			if w.localGoroutineQueues[threadIdx].IsEmpty() {
				continue
			}
			// pop from bottom of self queue
			currTask, _ := (w.localGoroutineQueues[threadIdx].PopBottom()).(Request)
			processImage(currTask, currTask.dataDir)
			w.mtx.Lock()
			w.totalTasks -= 1
			w.mtx.Unlock()
		}

		w.mtx.Lock()
		if w.localGoroutineQueues[threadIdx].IsEmpty() && *w.shutdown && w.totalTasks == 0 {
			w.mtx.Unlock()
			break
		}
		w.mtx.Unlock()
	}

}

func (w *WorkStealingExecutor) Shutdown() {
	w.mtx.Lock()
	*w.shutdown = true
	w.mtx.Unlock()
	w.wg.Wait()
}

func (w *WorkStealingExecutor) Submit(task interface{}) {
	w.mtx.Lock()
	w.localGoroutineQueues[w.currIndex].PushBottom(task)
	w.totalTasks += 1
	w.currIndex += 1
	w.currIndex = w.currIndex % w.capacity
	w.mtx.Unlock()
}