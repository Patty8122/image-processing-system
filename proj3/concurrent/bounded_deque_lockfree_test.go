package concurrent

import (
	"sync"
	"testing"
	// "fmt"
)

func TestBoundedDEQueue(t *testing.T) {
	const capacity = 500
	queue := NewBoundedDEQueue(capacity)

	// Test basic push and pop operations
	queue.PushBottom("Task1")
	queue.PushBottom("Task2")

	if size := queue.Size(); size != 2 {
		t.Errorf("Expected size %d, got %d", 2, size)
	}

	task := queue.PopTop()
	if task != "Task1" {
		t.Errorf("Expected popped task %s, got %v", "Task1", task)
	}

	// Test concurrent push and pop operations
	var wg sync.WaitGroup
	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			task := Task(taskID)
			fmt.Println(queue.PushBottom(task))
		}(i)
	}

	wg.Wait()

	// The size should be equal to the number of goroutines (numGoroutines)
	if size := queue.Size(); size != numGoroutines+1 { // +1 for the initial push
		t.Errorf("Expected size %d, got %d", numGoroutines+1, size)
	}

	// Test concurrent pop operations from both ends
	var poppedTasksTop, poppedTasksBottom []Task
	var poppedTasksMutex sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskTop := queue.PopTop()
			taskBottom := queue.PopBottom()
			poppedTasksMutex.Lock()
			if taskTop != nil {
				poppedTasksTop = append(poppedTasksTop, taskTop)
			}
			if taskBottom != nil {
				poppedTasksBottom = append(poppedTasksBottom, taskBottom)
			}
			poppedTasksMutex.Unlock()
		}()
	}

	wg.Wait()

	// The poppedTasksTop and poppedTasksBottom slices should each contain numGoroutines tasks
	if len(poppedTasksTop) != numGoroutines {
		t.Errorf("Expected %d popped tasks from top, got %d", numGoroutines, len(poppedTasksTop))
	}
	if len(poppedTasksBottom) != numGoroutines {
		t.Errorf("Expected %d popped tasks from bottom, got %d", numGoroutines, len(poppedTasksBottom))
	}

	// Ensure that all popped tasks are unique
	seenTasks := make(map[Task]struct{})
	for _, task := range poppedTasksTop {
		if _, exists := seenTasks[task]; exists {
			t.Errorf("Duplicate popped task from top: %v", task)
		}
		seenTasks[task] = struct{}{}
	}
	for _, task := range poppedTasksBottom {
		if _, exists := seenTasks[task]; exists {
			t.Errorf("Duplicate popped task from bottom: %v", task)
		}
		seenTasks[task] = struct{}{}
	}

	// Ensure the size is correct after all pops
	if size := queue.Size(); size != 0 {
		t.Errorf("Expected size %d, got %d", 0, size)
	}
}


func TestPopTopConcurrent(t *testing.T) {
	const capacity = 500
	queue := NewBoundedDEQueue(capacity)

	// Push tasks to the queue
	for i := 0; i < capacity; i++ {
		queue.PushBottom(i)
	}

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Test concurrent PopTop operations
	var poppedTasksTop []Task
	var poppedTasksMutex sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskTop := queue.PopTop()
			poppedTasksMutex.Lock()
			if taskTop != nil {
				poppedTasksTop = append(poppedTasksTop, taskTop)
			}
			poppedTasksMutex.Unlock()
		}()
	}

	wg.Wait()

	// The poppedTasksTop slice should contain numGoroutines tasks
	if len(poppedTasksTop) != numGoroutines {
		t.Errorf("Expected %d popped tasks from top, got %d", numGoroutines, len(poppedTasksTop))
	}

	// Ensure that all popped tasks are unique
	seenTasks := make(map[Task]struct{})
	for _, task := range poppedTasksTop {
		if _, exists := seenTasks[task]; exists {
			t.Errorf("Duplicate popped task from top: %v", task)
		}
		seenTasks[task] = struct{}{}
	}

	// Ensure the size is correct after all pops
	if size := queue.Size(); size != capacity-numGoroutines {
		t.Errorf("Expected size %d, got %d", capacity-numGoroutines, size)
	}
}


func TestPopTopAndPopBottomConcurrent(t *testing.T) {
	const capacity = 4
	queue := NewBoundedDEQueue(capacity)

	// Push tasks to the queue
	for i := 0; i < capacity; i++ {
		queue.PushBottom(i)
	}

	var wg sync.WaitGroup
	const numGoroutines = 2

	// Test concurrent PopTop and PopBottom operations
	var poppedTasksTop, poppedTasksBottom []Task
	var poppedTasksMutex sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskTop := queue.PopTop()
			taskBottom := queue.PopBottom()

			poppedTasksMutex.Lock()
			if taskTop != nil {
				poppedTasksTop = append(poppedTasksTop, taskTop)
			}
			if taskBottom != nil {
				poppedTasksBottom = append(poppedTasksBottom, taskBottom)
			}
			poppedTasksMutex.Unlock()
		}()
	}

	wg.Wait()

	// The poppedTasksTop and poppedTasksBottom slices should each contain numGoroutines tasks
	if len(poppedTasksTop) != numGoroutines {
		t.Errorf("Expected %d popped tasks from top, got %d", numGoroutines, len(poppedTasksTop))
	}
	if len(poppedTasksBottom) != numGoroutines {
		t.Errorf("Expected %d popped tasks from bottom, got %d", numGoroutines, len(poppedTasksBottom))
	}

	// Ensure that all popped tasks are unique
	seenTasks := make(map[Task]struct{})
	for _, task := range poppedTasksTop {
		if _, exists := seenTasks[task]; exists {
			t.Errorf("Duplicate popped task from top: %v", task)
		}
		seenTasks[task] = struct{}{}
	}
	for _, task := range poppedTasksBottom {
		if _, exists := seenTasks[task]; exists {
			t.Errorf("Duplicate popped task from bottom: %v", task)
		}
		seenTasks[task] = struct{}{}
	}

	// Ensure the size is correct after all pops
	if size := queue.Size(); size != capacity-numGoroutines*2 {
		t.Errorf("Expected size %d, got %d", capacity-numGoroutines*2, size)
	}
}

func TestMain(m *testing.M) {
	m.Run()
}