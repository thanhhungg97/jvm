package runtime

import (
	"fmt"
	"sync"
)

// CallbackExecutor provides the ability to invoke Java methods from Go
// This enables re-entrant execution: Go -> Java -> Go -> Java
type CallbackExecutor interface {
	// InvokeMethod invokes a method on an object
	// For Runnable, this calls run() with descriptor ()V
	InvokeMethod(obj interface{}, methodName, descriptor string) error

	// InvokeRunnable is a convenience method for Runnable.run()
	InvokeRunnable(runnable interface{}) error
}

// Global callback executor - set by the interpreter
var callbackExecutor CallbackExecutor
var callbackMu sync.RWMutex

// SetCallbackExecutor sets the global callback executor
func SetCallbackExecutor(executor CallbackExecutor) {
	callbackMu.Lock()
	defer callbackMu.Unlock()
	callbackExecutor = executor
}

// GetCallbackExecutor returns the global callback executor
func GetCallbackExecutor() CallbackExecutor {
	callbackMu.RLock()
	defer callbackMu.RUnlock()
	return callbackExecutor
}

// InvokeRunnable invokes the run() method on a Runnable object
func InvokeRunnable(runnable interface{}) error {
	executor := GetCallbackExecutor()
	if executor == nil {
		return fmt.Errorf("no callback executor registered")
	}
	return executor.InvokeRunnable(runnable)
}

// InvokeMethod invokes a method on an object
func InvokeMethod(obj interface{}, methodName, descriptor string) error {
	executor := GetCallbackExecutor()
	if executor == nil {
		return fmt.Errorf("no callback executor registered")
	}
	return executor.InvokeMethod(obj, methodName, descriptor)
}

// RunnableTask wraps a Runnable object for deferred execution
type RunnableTask struct {
	Runnable interface{} // The Java Runnable object
	Name     string      // Task name for debugging
}

// Execute runs the Runnable's run() method
func (rt *RunnableTask) Execute() error {
	return InvokeRunnable(rt.Runnable)
}

// PendingCallback represents a callback waiting to be executed
type PendingCallback struct {
	Object     interface{}
	MethodName string
	Descriptor string
	Args       []interface{}
}

// CallbackQueue holds callbacks waiting to be executed
type CallbackQueue struct {
	callbacks []*PendingCallback
	mu        sync.Mutex
}

// NewCallbackQueue creates a new callback queue
func NewCallbackQueue() *CallbackQueue {
	return &CallbackQueue{
		callbacks: make([]*PendingCallback, 0),
	}
}

// Enqueue adds a callback to the queue
func (cq *CallbackQueue) Enqueue(cb *PendingCallback) {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	cq.callbacks = append(cq.callbacks, cb)
}

// Dequeue removes and returns the next callback
func (cq *CallbackQueue) Dequeue() *PendingCallback {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	if len(cq.callbacks) == 0 {
		return nil
	}
	cb := cq.callbacks[0]
	cq.callbacks = cq.callbacks[1:]
	return cb
}

// Len returns the number of pending callbacks
func (cq *CallbackQueue) Len() int {
	cq.mu.Lock()
	defer cq.mu.Unlock()
	return len(cq.callbacks)
}


