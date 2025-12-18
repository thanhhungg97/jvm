package runtime

import (
	"fmt"
	"sync"
)

// Task execution registry - maps task IDs to their execution function
var taskRegistry = struct {
	sync.RWMutex
	tasks map[int32]func()
}{
	tasks: make(map[int32]func()),
}

// RegisterTask registers a task callback
func RegisterTask(id int32, callback func()) {
	taskRegistry.Lock()
	taskRegistry.tasks[id] = callback
	taskRegistry.Unlock()
}

// GetTask retrieves a registered task
func GetTask(id int32) func() {
	taskRegistry.RLock()
	defer taskRegistry.RUnlock()
	return taskRegistry.tasks[id]
}

func init() {
	// Register event loop native methods
	Natives.Register("EventLoop", "submit", "(ILjava/lang/String;)V", nativeEventLoopSubmit)
	Natives.Register("EventLoop", "submitRunnable", "(Ljava/lang/Runnable;)V", nativeEventLoopSubmitRunnable)
	Natives.Register("EventLoop", "setTimeout", "(ILjava/lang/String;J)V", nativeEventLoopSetTimeout)
	Natives.Register("EventLoop", "setTimeoutRunnable", "(Ljava/lang/Runnable;J)V", nativeEventLoopSetTimeoutRunnable)
	Natives.Register("EventLoop", "setInterval", "(ILjava/lang/String;J)V", nativeEventLoopSetInterval)
	Natives.Register("EventLoop", "run", "()V", nativeEventLoopRun)
	Natives.Register("EventLoop", "stop", "()V", nativeEventLoopStop)
	Natives.Register("EventLoop", "isRunning", "()Z", nativeEventLoopIsRunning)
	Natives.Register("EventLoop", "printStats", "()V", nativeEventLoopPrintStats)
	Natives.Register("EventLoop", "reset", "()V", nativeEventLoopReset)

	// Also register as EventLoopDemo for the demo class
	Natives.Register("EventLoopDemo", "submit", "(ILjava/lang/String;)V", nativeEventLoopSubmit)
	Natives.Register("EventLoopDemo", "submitRunnable", "(Ljava/lang/Runnable;)V", nativeEventLoopSubmitRunnable)
	Natives.Register("EventLoopDemo", "setTimeout", "(ILjava/lang/String;J)V", nativeEventLoopSetTimeout)
	Natives.Register("EventLoopDemo", "setTimeoutRunnable", "(Ljava/lang/Runnable;J)V", nativeEventLoopSetTimeoutRunnable)
	Natives.Register("EventLoopDemo", "setInterval", "(ILjava/lang/String;J)V", nativeEventLoopSetInterval)
	Natives.Register("EventLoopDemo", "run", "()V", nativeEventLoopRun)
	Natives.Register("EventLoopDemo", "stop", "()V", nativeEventLoopStop)
	Natives.Register("EventLoopDemo", "isRunning", "()Z", nativeEventLoopIsRunning)
	Natives.Register("EventLoopDemo", "printStats", "()V", nativeEventLoopPrintStats)
	Natives.Register("EventLoopDemo", "reset", "()V", nativeEventLoopReset)
}

// nativeEventLoopSubmit submits a task to the event loop
// Java: static native void submit(int taskId, String name)
func nativeEventLoopSubmit(frame *Frame) error {
	stack := frame.OperandStack
	nameRef := stack.PopRef()
	taskId := stack.PopInt()

	name := "task"
	if s, ok := nameRef.(string); ok {
		name = s
	}

	el := GetEventLoop()

	// Create a callback that prints the task execution
	callback := func() {
		fmt.Printf("[%s] Task %d executing\n", name, taskId)
	}

	el.Submit(taskId, name, callback)
	return nil
}

// nativeEventLoopSubmitRunnable submits a Runnable to the event loop
// Java: static native void submitRunnable(Runnable r)
func nativeEventLoopSubmitRunnable(frame *Frame) error {
	stack := frame.OperandStack
	runnableRef := stack.PopRef()

	if runnableRef == nil {
		return fmt.Errorf("NullPointerException: runnable is null")
	}

	el := GetEventLoop()
	taskId := int32(el.taskCount + 1)

	// Create a callback that invokes the Runnable's run() method
	callback := func() {
		if err := InvokeRunnable(runnableRef); err != nil {
			fmt.Printf("[lambda] Error executing runnable: %v\n", err)
		}
	}

	el.Submit(taskId, "lambda", callback)
	return nil
}

// nativeEventLoopSetTimeout schedules a delayed task
// Java: static native void setTimeout(int taskId, String name, long delayMs)
func nativeEventLoopSetTimeout(frame *Frame) error {
	stack := frame.OperandStack
	delayMs := stack.PopLong()
	nameRef := stack.PopRef()
	taskId := stack.PopInt()

	name := "timer"
	if s, ok := nameRef.(string); ok {
		name = s
	}

	el := GetEventLoop()

	callback := func() {
		fmt.Printf("[%s] Timer %d fired after %dms\n", name, taskId, delayMs)
	}

	el.SetTimeout(taskId, name, delayMs, callback)
	return nil
}

// nativeEventLoopSetTimeoutRunnable schedules a Runnable with delay
// Java: static native void setTimeoutRunnable(Runnable r, long delayMs)
func nativeEventLoopSetTimeoutRunnable(frame *Frame) error {
	stack := frame.OperandStack
	delayMs := stack.PopLong()
	runnableRef := stack.PopRef()

	if runnableRef == nil {
		return fmt.Errorf("NullPointerException: runnable is null")
	}

	el := GetEventLoop()
	taskId := int32(el.timerCount + 100)

	callback := func() {
		if err := InvokeRunnable(runnableRef); err != nil {
			fmt.Printf("[timer-lambda] Error executing runnable: %v\n", err)
		}
	}

	el.SetTimeout(taskId, "timer-lambda", delayMs, callback)
	return nil
}

// nativeEventLoopSetInterval schedules a repeating task
// Java: static native void setInterval(int taskId, String name, long periodMs)
func nativeEventLoopSetInterval(frame *Frame) error {
	stack := frame.OperandStack
	periodMs := stack.PopLong()
	nameRef := stack.PopRef()
	taskId := stack.PopInt()

	name := "interval"
	if s, ok := nameRef.(string); ok {
		name = s
	}

	el := GetEventLoop()

	counter := 0
	callback := func() {
		counter++
		fmt.Printf("[%s] Interval %d tick #%d\n", name, taskId, counter)
	}

	el.SetInterval(taskId, name, periodMs, callback)
	return nil
}

// nativeEventLoopRun starts the event loop
// Java: static native void run()
func nativeEventLoopRun(frame *Frame) error {
	el := GetEventLoop()
	fmt.Println("Event loop starting...")
	el.Run()
	fmt.Println("Event loop finished.")
	return nil
}

// nativeEventLoopStop stops the event loop
// Java: static native void stop()
func nativeEventLoopStop(frame *Frame) error {
	el := GetEventLoop()
	el.Stop()
	return nil
}

// nativeEventLoopIsRunning checks if event loop is running
// Java: static native boolean isRunning()
func nativeEventLoopIsRunning(frame *Frame) error {
	el := GetEventLoop()
	if el.IsRunning() {
		frame.OperandStack.PushInt(1)
	} else {
		frame.OperandStack.PushInt(0)
	}
	return nil
}

// nativeEventLoopPrintStats prints event loop statistics
// Java: static native void printStats()
func nativeEventLoopPrintStats(frame *Frame) error {
	el := GetEventLoop()
	el.PrintStats()
	return nil
}

// nativeEventLoopReset resets the event loop
// Java: static native void reset()
func nativeEventLoopReset(frame *Frame) error {
	ResetEventLoop()
	return nil
}
