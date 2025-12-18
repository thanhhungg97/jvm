package runtime

import (
	"container/heap"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Scheduler manages green thread (fiber) execution
type Scheduler struct {
	jvm          *JVM
	fibers       map[int64]*Fiber
	readyQueue   *FiberQueue
	currentFiber *Fiber
	mu           sync.Mutex
	running      bool
	workersCount int
	workersDone  sync.WaitGroup
	tickInterval time.Duration
	stats        SchedulerStats
	onFiberDone  func(*Fiber)
}

// SchedulerStats tracks scheduler performance
type SchedulerStats struct {
	FibersCreated   int64
	FibersCompleted int64
	ContextSwitches int64
	TotalYields     int64
}

// NewScheduler creates a new fiber scheduler
func NewScheduler(jvm *JVM) *Scheduler {
	return &Scheduler{
		jvm:          jvm,
		fibers:       make(map[int64]*Fiber),
		readyQueue:   NewFiberQueue(),
		running:      false,
		workersCount: 1, // Single worker for deterministic execution
		tickInterval: time.Microsecond * 100,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	// Start worker goroutines
	for i := 0; i < s.workersCount; i++ {
		s.workersDone.Add(1)
		go s.worker(i)
	}
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	s.workersDone.Wait()
}

// IsRunning returns true if scheduler is running
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// Spawn creates and schedules a new fiber
func (s *Scheduler) Spawn(name string, task func(*Fiber)) *Fiber {
	thread := NewThread()
	if s.jvm != nil {
		thread = s.jvm.CreateThread()
	}

	fiber := NewFiber(name, thread)
	fiber.scheduler = s

	s.mu.Lock()
	s.fibers[fiber.ID] = fiber
	atomic.AddInt64(&s.stats.FibersCreated, 1)
	s.mu.Unlock()

	// Wrap the task to handle completion
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fiber.err = fmt.Errorf("fiber panic: %v", r)
			}
			fiber.SetState(FiberDead)
			atomic.AddInt64(&s.stats.FibersCompleted, 1)
			if s.onFiberDone != nil {
				s.onFiberDone(fiber)
			}
		}()

		fiber.SetState(FiberRunning)
		task(fiber)
	}()

	return fiber
}

// SpawnMethod spawns a fiber that executes a JVM method
func (s *Scheduler) SpawnMethod(name string, executor func(*Fiber) error) *Fiber {
	return s.Spawn(name, func(f *Fiber) {
		if err := executor(f); err != nil {
			f.err = err
		}
	})
}

// Yield causes the current fiber to yield execution
func (s *Scheduler) Yield(fiber *Fiber) {
	atomic.AddInt64(&s.stats.TotalYields, 1)
	// In our goroutine-based model, we use runtime.Gosched equivalent
	time.Sleep(time.Microsecond)
}

// Sleep puts a fiber to sleep for the specified duration
func (s *Scheduler) Sleep(fiber *Fiber, millis int64) {
	fiber.SetState(FiberSleeping)
	time.Sleep(time.Duration(millis) * time.Millisecond)
	fiber.SetState(FiberRunning)
}

// worker is a scheduler worker goroutine
func (s *Scheduler) worker(id int) {
	defer s.workersDone.Done()

	for {
		s.mu.Lock()
		if !s.running {
			s.mu.Unlock()
			return
		}
		s.mu.Unlock()

		time.Sleep(s.tickInterval)
	}
}

// GetStats returns scheduler statistics
func (s *Scheduler) GetStats() SchedulerStats {
	return SchedulerStats{
		FibersCreated:   atomic.LoadInt64(&s.stats.FibersCreated),
		FibersCompleted: atomic.LoadInt64(&s.stats.FibersCompleted),
		ContextSwitches: atomic.LoadInt64(&s.stats.ContextSwitches),
		TotalYields:     atomic.LoadInt64(&s.stats.TotalYields),
	}
}

// FiberCount returns the number of active fibers
func (s *Scheduler) FiberCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, f := range s.fibers {
		if f.IsAlive() {
			count++
		}
	}
	return count
}

// WaitAll waits for all fibers to complete
func (s *Scheduler) WaitAll() {
	for s.FiberCount() > 0 {
		time.Sleep(time.Millisecond)
	}
}

// WaitFor waits for a specific fiber to complete
func (s *Scheduler) WaitFor(fiber *Fiber) {
	for fiber.IsAlive() {
		time.Sleep(time.Millisecond)
	}
}

// OnFiberDone sets a callback for when fibers complete
func (s *Scheduler) OnFiberDone(callback func(*Fiber)) {
	s.onFiberDone = callback
}

// PrintStats prints scheduler statistics
func (s *Scheduler) PrintStats() {
	stats := s.GetStats()
	fmt.Printf("=== Scheduler Stats ===\n")
	fmt.Printf("Fibers Created:   %d\n", stats.FibersCreated)
	fmt.Printf("Fibers Completed: %d\n", stats.FibersCompleted)
	fmt.Printf("Context Switches: %d\n", stats.ContextSwitches)
	fmt.Printf("Total Yields:     %d\n", stats.TotalYields)
}

// =============== Priority Queue for Fiber Scheduling ===============

// FiberQueue is a priority queue for fibers
type FiberQueue struct {
	fibers []*Fiber
	mu     sync.Mutex
}

// NewFiberQueue creates a new fiber queue
func NewFiberQueue() *FiberQueue {
	fq := &FiberQueue{
		fibers: make([]*Fiber, 0),
	}
	heap.Init(fq)
	return fq
}

// Len returns the queue length
func (fq *FiberQueue) Len() int {
	return len(fq.fibers)
}

// Less compares priorities (higher priority first)
func (fq *FiberQueue) Less(i, j int) bool {
	return fq.fibers[i].priority > fq.fibers[j].priority
}

// Swap swaps two elements
func (fq *FiberQueue) Swap(i, j int) {
	fq.fibers[i], fq.fibers[j] = fq.fibers[j], fq.fibers[i]
}

// Push adds a fiber to the queue
func (fq *FiberQueue) Push(x interface{}) {
	fq.fibers = append(fq.fibers, x.(*Fiber))
}

// Pop removes and returns the highest priority fiber
func (fq *FiberQueue) Pop() interface{} {
	n := len(fq.fibers)
	fiber := fq.fibers[n-1]
	fq.fibers = fq.fibers[:n-1]
	return fiber
}

// Enqueue adds a fiber to the queue (thread-safe)
func (fq *FiberQueue) Enqueue(fiber *Fiber) {
	fq.mu.Lock()
	defer fq.mu.Unlock()
	heap.Push(fq, fiber)
}

// Dequeue removes and returns the highest priority fiber (thread-safe)
func (fq *FiberQueue) Dequeue() *Fiber {
	fq.mu.Lock()
	defer fq.mu.Unlock()
	if len(fq.fibers) == 0 {
		return nil
	}
	return heap.Pop(fq).(*Fiber)
}

// IsEmpty returns true if queue is empty
func (fq *FiberQueue) IsEmpty() bool {
	fq.mu.Lock()
	defer fq.mu.Unlock()
	return len(fq.fibers) == 0
}
