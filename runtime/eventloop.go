package runtime

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// Task represents a unit of work in the event loop
type Task struct {
	ID       int32
	Name     string
	Callback func()
}

// TimerTask represents a scheduled task with a deadline
type TimerTask struct {
	ID       int32
	Name     string
	Deadline time.Time
	Callback func()
	Interval time.Duration // For setInterval, 0 for setTimeout
	index    int           // Index in the heap
}

// EventLoop implements a Node.js-style event loop
type EventLoop struct {
	tasks      chan Task     // Task queue (FIFO)
	timers     *TimerHeap    // Timer heap (min-heap by deadline)
	running    bool          // Is the loop running?
	stopped    chan struct{} // Signal to stop
	mu         sync.Mutex
	taskCount  int32 // Number of tasks processed
	timerCount int32 // Number of timers fired
}

// Global event loop instance
var globalEventLoop *EventLoop
var eventLoopOnce sync.Once

// GetEventLoop returns the global event loop, creating it if needed
func GetEventLoop() *EventLoop {
	eventLoopOnce.Do(func() {
		globalEventLoop = NewEventLoop()
	})
	return globalEventLoop
}

// ResetEventLoop resets the global event loop (for testing)
func ResetEventLoop() {
	globalEventLoop = NewEventLoop()
}

// NewEventLoop creates a new event loop
func NewEventLoop() *EventLoop {
	return &EventLoop{
		tasks:   make(chan Task, 1000), // Buffered channel for tasks
		timers:  NewTimerHeap(),
		stopped: make(chan struct{}),
	}
}

// Submit adds a task to the queue
func (el *EventLoop) Submit(id int32, name string, callback func()) {
	el.tasks <- Task{
		ID:       id,
		Name:     name,
		Callback: callback,
	}
}

// SetTimeout schedules a task to run after a delay
func (el *EventLoop) SetTimeout(id int32, name string, delayMs int64, callback func()) {
	el.mu.Lock()
	defer el.mu.Unlock()

	timer := &TimerTask{
		ID:       id,
		Name:     name,
		Deadline: time.Now().Add(time.Duration(delayMs) * time.Millisecond),
		Callback: callback,
		Interval: 0,
	}
	heap.Push(el.timers, timer)
}

// SetInterval schedules a repeating task
func (el *EventLoop) SetInterval(id int32, name string, periodMs int64, callback func()) {
	el.mu.Lock()
	defer el.mu.Unlock()

	timer := &TimerTask{
		ID:       id,
		Name:     name,
		Deadline: time.Now().Add(time.Duration(periodMs) * time.Millisecond),
		Callback: callback,
		Interval: time.Duration(periodMs) * time.Millisecond,
	}
	heap.Push(el.timers, timer)
}

// Run starts the event loop and blocks until stopped or all tasks complete
func (el *EventLoop) Run() {
	el.mu.Lock()
	if el.running {
		el.mu.Unlock()
		return
	}
	el.running = true
	el.stopped = make(chan struct{})
	el.mu.Unlock()

	for {
		// Check if we should stop
		select {
		case <-el.stopped:
			return
		default:
		}

		// 1. Fire any ready timers
		el.fireReadyTimers()

		// 2. Process one task from the queue (non-blocking)
		select {
		case task := <-el.tasks:
			el.taskCount++
			if task.Callback != nil {
				task.Callback()
			}
		case <-el.stopped:
			return
		default:
			// No task available, check if we're done
			el.mu.Lock()
			hasTimers := el.timers.Len() > 0
			el.mu.Unlock()

			if !hasTimers && len(el.tasks) == 0 {
				// No more work to do
				el.mu.Lock()
				el.running = false
				el.mu.Unlock()
				return
			}

			// Sleep a bit before checking again
			time.Sleep(time.Millisecond)
		}
	}
}

// RunFor runs the event loop for a maximum duration
func (el *EventLoop) RunFor(maxDuration time.Duration) {
	done := make(chan struct{})
	go func() {
		el.Run()
		close(done)
	}()

	select {
	case <-done:
		// Loop finished naturally
	case <-time.After(maxDuration):
		el.Stop()
	}
}

// Stop stops the event loop
func (el *EventLoop) Stop() {
	el.mu.Lock()
	defer el.mu.Unlock()

	if el.running {
		el.running = false
		close(el.stopped)
	}
}

// IsRunning returns true if the event loop is running
func (el *EventLoop) IsRunning() bool {
	el.mu.Lock()
	defer el.mu.Unlock()
	return el.running
}

// fireReadyTimers fires all timers whose deadline has passed
func (el *EventLoop) fireReadyTimers() {
	now := time.Now()

	el.mu.Lock()
	defer el.mu.Unlock()

	for el.timers.Len() > 0 {
		// Peek at the next timer
		timer := (*el.timers)[0]
		if timer.Deadline.After(now) {
			// Not ready yet
			break
		}

		// Pop and fire the timer
		heap.Pop(el.timers)
		el.timerCount++

		if timer.Callback != nil {
			timer.Callback()
		}

		// If it's an interval, reschedule it
		if timer.Interval > 0 {
			timer.Deadline = now.Add(timer.Interval)
			heap.Push(el.timers, timer)
		}
	}
}

// Stats returns event loop statistics
func (el *EventLoop) Stats() (tasks int32, timers int32) {
	return el.taskCount, el.timerCount
}

// PendingTasks returns the number of pending tasks
func (el *EventLoop) PendingTasks() int {
	return len(el.tasks)
}

// PendingTimers returns the number of pending timers
func (el *EventLoop) PendingTimers() int {
	el.mu.Lock()
	defer el.mu.Unlock()
	return el.timers.Len()
}

// =============== Timer Heap (Min-Heap by Deadline) ===============

// TimerHeap is a min-heap of timer tasks ordered by deadline
type TimerHeap []*TimerTask

// NewTimerHeap creates a new timer heap
func NewTimerHeap() *TimerHeap {
	h := &TimerHeap{}
	heap.Init(h)
	return h
}

func (h TimerHeap) Len() int           { return len(h) }
func (h TimerHeap) Less(i, j int) bool { return h[i].Deadline.Before(h[j].Deadline) }
func (h TimerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *TimerHeap) Push(x interface{}) {
	n := len(*h)
	timer := x.(*TimerTask)
	timer.index = n
	*h = append(*h, timer)
}

func (h *TimerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	timer := old[n-1]
	old[n-1] = nil   // Avoid memory leak
	timer.index = -1 // Mark as removed
	*h = old[0 : n-1]
	return timer
}

// PrintStats prints event loop statistics
func (el *EventLoop) PrintStats() {
	tasks, timers := el.Stats()
	fmt.Println("=== Event Loop Statistics ===")
	fmt.Printf("Tasks Processed:  %d\n", tasks)
	fmt.Printf("Timers Fired:     %d\n", timers)
	fmt.Printf("Pending Tasks:    %d\n", el.PendingTasks())
	fmt.Printf("Pending Timers:   %d\n", el.PendingTimers())
}
