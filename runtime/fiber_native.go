package runtime

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Global scheduler for fibers (reserved for future use)
var _ *Scheduler // globalScheduler - reserved
var fiberOutputMu sync.Mutex

// Channel for fiber task execution
type FiberTask struct {
	ID       int64
	Name     string
	TaskID   int32 // The task ID passed from Java
	Complete chan struct{}
	Result   int32
}

// Active fibers storage
var activeFibers = struct {
	sync.RWMutex
	fibers map[int64]*FiberTask
}{
	fibers: make(map[int64]*FiberTask),
}

var fiberTaskCounter int64

func init() {
	// Register fiber/green thread natives
	Natives.Register("Fiber", "spawn", "(ILjava/lang/String;)J", nativeFiberSpawn)
	Natives.Register("Fiber", "yield", "()V", nativeFiberYield)
	Natives.Register("Fiber", "sleep", "(J)V", nativeFiberSleep)
	Natives.Register("Fiber", "join", "(J)V", nativeFiberJoin)
	Natives.Register("Fiber", "isAlive", "(J)Z", nativeFiberIsAlive)
	Natives.Register("Fiber", "current", "()J", nativeFiberCurrent)
	Natives.Register("Fiber", "count", "()I", nativeFiberCount)
	Natives.Register("Fiber", "printStats", "()V", nativeFiberPrintStats)

	// Also register with GreenThreads class name (Java file uses plural)
	Natives.Register("GreenThreads", "spawn", "(ILjava/lang/String;)J", nativeFiberSpawn)
	Natives.Register("GreenThreads", "yield", "()V", nativeFiberYield)
	Natives.Register("GreenThreads", "sleep", "(J)V", nativeFiberSleep)
	Natives.Register("GreenThreads", "join", "(J)V", nativeFiberJoin)
	Natives.Register("GreenThreads", "isAlive", "(J)Z", nativeFiberIsAlive)
	Natives.Register("GreenThreads", "current", "()J", nativeFiberCurrent)
	Natives.Register("GreenThreads", "count", "()I", nativeFiberCount)
	Natives.Register("GreenThreads", "printStats", "()V", nativeFiberPrintStats)

	// Parallel execution helpers
	Natives.Register("Parallel", "run", "(I)V", nativeParallelRun)
	Natives.Register("Parallel", "forEach", "(II)V", nativeParallelForEach)
}

// nativeFiberSpawn spawns a new fiber
// Java signature: static native long spawn(int taskId, String name)
func nativeFiberSpawn(frame *Frame) error {
	stack := frame.OperandStack
	nameRef := stack.PopRef()
	taskID := stack.PopInt()

	name := "fiber"
	if s, ok := nameRef.(string); ok {
		name = s
	}

	fiberID := atomic.AddInt64(&fiberTaskCounter, 1)

	task := &FiberTask{
		ID:       fiberID,
		Name:     name,
		TaskID:   taskID,
		Complete: make(chan struct{}),
	}

	activeFibers.Lock()
	activeFibers.fibers[fiberID] = task
	activeFibers.Unlock()

	// Spawn a goroutine to execute the task
	go func() {
		defer func() {
			close(task.Complete)
		}()

		// Simulate work based on taskID
		// In a full implementation, this would call back into the JVM
		// to execute a Java method
		executeTask(task)
	}()

	stack.PushLong(fiberID)
	return nil
}

// executeTask simulates executing a fiber task
func executeTask(task *FiberTask) {
	// This simulates the fiber doing work
	// Each task does some iterations with yields
	iterations := int(task.TaskID) * 3

	for i := 0; i < iterations; i++ {
		// Simulate some work
		time.Sleep(time.Millisecond * 10)

		// Print progress (thread-safe)
		fiberOutputMu.Lock()
		fmt.Printf("[%s] iteration %d/%d\n", task.Name, i+1, iterations)
		fiberOutputMu.Unlock()

		// Yield to other fibers
		time.Sleep(time.Microsecond * 100)
	}

	task.Result = task.TaskID * 10 // Result is taskID * 10
}

// nativeFiberYield yields execution to other fibers
func nativeFiberYield(frame *Frame) error {
	time.Sleep(time.Microsecond * 100)
	return nil
}

// nativeFiberSleep puts the fiber to sleep
func nativeFiberSleep(frame *Frame) error {
	millis := frame.OperandStack.PopLong()
	time.Sleep(time.Duration(millis) * time.Millisecond)
	return nil
}

// nativeFiberJoin waits for a fiber to complete
func nativeFiberJoin(frame *Frame) error {
	fiberID := frame.OperandStack.PopLong()

	activeFibers.RLock()
	task, exists := activeFibers.fibers[fiberID]
	activeFibers.RUnlock()

	if exists {
		<-task.Complete
	}

	return nil
}

// nativeFiberIsAlive checks if a fiber is still running
func nativeFiberIsAlive(frame *Frame) error {
	fiberID := frame.OperandStack.PopLong()

	activeFibers.RLock()
	task, exists := activeFibers.fibers[fiberID]
	activeFibers.RUnlock()

	if !exists {
		frame.OperandStack.PushInt(0)
		return nil
	}

	select {
	case <-task.Complete:
		frame.OperandStack.PushInt(0)
	default:
		frame.OperandStack.PushInt(1)
	}

	return nil
}

// nativeFiberCurrent returns the current fiber ID (or 0 for main)
func nativeFiberCurrent(frame *Frame) error {
	// In this simple implementation, we return 0 for main thread
	frame.OperandStack.PushLong(0)
	return nil
}

// nativeFiberCount returns the number of active fibers
func nativeFiberCount(frame *Frame) error {
	activeFibers.RLock()
	count := 0
	for _, task := range activeFibers.fibers {
		select {
		case <-task.Complete:
			// Completed
		default:
			count++
		}
	}
	activeFibers.RUnlock()

	frame.OperandStack.PushInt(int32(count))
	return nil
}

// nativeFiberPrintStats prints fiber statistics
func nativeFiberPrintStats(frame *Frame) error {
	activeFibers.RLock()
	total := len(activeFibers.fibers)
	active := 0
	completed := 0
	for _, task := range activeFibers.fibers {
		select {
		case <-task.Complete:
			completed++
		default:
			active++
		}
	}
	activeFibers.RUnlock()

	fmt.Println("=== Fiber Statistics ===")
	fmt.Printf("Total Created: %d\n", total)
	fmt.Printf("Active:        %d\n", active)
	fmt.Printf("Completed:     %d\n", completed)

	return nil
}

// nativeParallelRun runs N tasks in parallel fibers
func nativeParallelRun(frame *Frame) error {
	numTasks := frame.OperandStack.PopInt()

	var wg sync.WaitGroup
	tasks := make([]*FiberTask, numTasks)

	// Spawn all tasks
	for i := int32(0); i < numTasks; i++ {
		fiberID := atomic.AddInt64(&fiberTaskCounter, 1)
		task := &FiberTask{
			ID:       fiberID,
			Name:     fmt.Sprintf("parallel-%d", i),
			TaskID:   i + 1,
			Complete: make(chan struct{}),
		}
		tasks[i] = task

		activeFibers.Lock()
		activeFibers.fibers[fiberID] = task
		activeFibers.Unlock()

		wg.Add(1)
		go func(t *FiberTask) {
			defer wg.Done()
			defer close(t.Complete)
			executeTask(t)
		}(task)
	}

	// Wait for all to complete
	wg.Wait()

	return nil
}

// nativeParallelForEach runs a parallel for-each over a range
func nativeParallelForEach(frame *Frame) error {
	end := frame.OperandStack.PopInt()
	start := frame.OperandStack.PopInt()

	var wg sync.WaitGroup

	for i := start; i < end; i++ {
		wg.Add(1)
		go func(idx int32) {
			defer wg.Done()
			fiberOutputMu.Lock()
			fmt.Printf("[parallel] processing index %d\n", idx)
			fiberOutputMu.Unlock()
			time.Sleep(time.Millisecond * 20)
		}(i)
	}

	wg.Wait()
	return nil
}
