package runtime

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// FiberState represents the state of a fiber
type FiberState int32

const (
	FiberReady    FiberState = iota // Ready to run
	FiberRunning                    // Currently executing
	FiberBlocked                    // Waiting for something
	FiberSleeping                   // Sleeping for a duration
	FiberDead                       // Finished execution
)

func (s FiberState) String() string {
	switch s {
	case FiberReady:
		return "READY"
	case FiberRunning:
		return "RUNNING"
	case FiberBlocked:
		return "BLOCKED"
	case FiberSleeping:
		return "SLEEPING"
	case FiberDead:
		return "DEAD"
	default:
		return "UNKNOWN"
	}
}

// Fiber represents a green thread (lightweight thread)
type Fiber struct {
	ID           int64
	Name         string
	State        FiberState
	Thread       *Thread       // JVM thread state (stack, frames)
	scheduler    *Scheduler    // Reference to scheduler
	wakeupChan   chan struct{} // Channel to wake up blocked fiber
	result       interface{}   // Result when fiber completes
	err          error         // Error if fiber failed
	priority     int           // Scheduling priority (higher = more priority)
	yieldCounter int64         // Number of instructions before yielding
	mu           sync.Mutex
}

// Global fiber ID counter
var fiberIDCounter int64

// NewFiber creates a new fiber
func NewFiber(name string, thread *Thread) *Fiber {
	id := atomic.AddInt64(&fiberIDCounter, 1)
	return &Fiber{
		ID:           id,
		Name:         name,
		State:        FiberReady,
		Thread:       thread,
		wakeupChan:   make(chan struct{}, 1),
		priority:     5,    // Default priority
		yieldCounter: 1000, // Yield after 1000 instructions
	}
}

// SetState atomically sets the fiber state
func (f *Fiber) SetState(state FiberState) {
	f.mu.Lock()
	f.State = state
	f.mu.Unlock()
}

// GetState atomically gets the fiber state
func (f *Fiber) GetState() FiberState {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.State
}

// IsAlive returns true if fiber hasn't completed
func (f *Fiber) IsAlive() bool {
	state := f.GetState()
	return state != FiberDead
}

// Yield voluntarily yields execution to other fibers
func (f *Fiber) Yield() {
	if f.scheduler != nil {
		f.scheduler.Yield(f)
	}
}

// Sleep puts the fiber to sleep for the given milliseconds
func (f *Fiber) Sleep(millis int64) {
	if f.scheduler != nil {
		f.scheduler.Sleep(f, millis)
	}
}

// Wake wakes up a sleeping or blocked fiber
func (f *Fiber) Wake() {
	select {
	case f.wakeupChan <- struct{}{}:
	default:
		// Already has a wake signal pending
	}
}

// String returns a string representation of the fiber
func (f *Fiber) String() string {
	return fmt.Sprintf("Fiber[%d:%s:%s]", f.ID, f.Name, f.State)
}

// FiberGroup manages a group of related fibers
type FiberGroup struct {
	Name    string
	fibers  []*Fiber
	mu      sync.Mutex
	done    chan struct{}
	results map[int64]interface{}
}

// NewFiberGroup creates a new fiber group
func NewFiberGroup(name string) *FiberGroup {
	return &FiberGroup{
		Name:    name,
		fibers:  make([]*Fiber, 0),
		done:    make(chan struct{}),
		results: make(map[int64]interface{}),
	}
}

// Add adds a fiber to the group
func (fg *FiberGroup) Add(fiber *Fiber) {
	fg.mu.Lock()
	defer fg.mu.Unlock()
	fg.fibers = append(fg.fibers, fiber)
}

// Size returns the number of fibers in the group
func (fg *FiberGroup) Size() int {
	fg.mu.Lock()
	defer fg.mu.Unlock()
	return len(fg.fibers)
}

// AllDone returns true if all fibers have completed
func (fg *FiberGroup) AllDone() bool {
	fg.mu.Lock()
	defer fg.mu.Unlock()
	for _, f := range fg.fibers {
		if f.IsAlive() {
			return false
		}
	}
	return true
}

// WaitAll waits for all fibers to complete
func (fg *FiberGroup) WaitAll() {
	for !fg.AllDone() {
		// Busy wait with yield
		// In a real implementation, we'd use a condition variable
	}
}

// GetResults returns results from all fibers
func (fg *FiberGroup) GetResults() map[int64]interface{} {
	fg.mu.Lock()
	defer fg.mu.Unlock()
	results := make(map[int64]interface{})
	for _, f := range fg.fibers {
		results[f.ID] = f.result
	}
	return results
}
