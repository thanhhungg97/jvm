package runtime

import (
	"fmt"
	"simplejvm/classfile"
	"sync"
	"sync/atomic"
)

// JVM represents the Java Virtual Machine instance
type JVM struct {
	// Thread management
	mainThread    *Thread
	threads       []*Thread
	threadCounter int64
	threadMutex   sync.RWMutex

	// Class loading
	classCache map[string]*classfile.ClassFile
	classMutex sync.RWMutex

	// Monitor management for synchronized blocks
	monitors     map[any]*Monitor
	monitorMutex sync.Mutex

	// Heap for object allocation (Phase 7)
	heap *Heap

	// Global state
	running atomic.Bool
}

// Monitor represents a Java monitor for synchronized blocks
type Monitor struct {
	owner      *Thread   // Current owner thread
	entryCount int       // Reentrant count
	waitSet    []*Thread // Threads waiting on this monitor
	mutex      sync.Mutex
	cond       *sync.Cond
}

// NewJVM creates a new JVM instance
func NewJVM() *JVM {
	jvm := &JVM{
		classCache: make(map[string]*classfile.ClassFile),
		monitors:   make(map[any]*Monitor),
		heap:       NewHeap(),
	}
	jvm.running.Store(true)
	return jvm
}

// CreateThread creates a new thread
func (jvm *JVM) CreateThread() *Thread {
	jvm.threadMutex.Lock()
	defer jvm.threadMutex.Unlock()

	id := atomic.AddInt64(&jvm.threadCounter, 1)
	thread := &Thread{
		id:      id,
		stack:   make([]*Frame, 0, 32),
		Classes: jvm.classCache,
		jvm:     jvm,
	}

	jvm.threads = append(jvm.threads, thread)
	if jvm.mainThread == nil {
		jvm.mainThread = thread
	}

	return thread
}

// GetMainThread returns the main thread
func (jvm *JVM) GetMainThread() *Thread {
	return jvm.mainThread
}

// LoadClass loads and caches a class
func (jvm *JVM) LoadClass(name string, cf *classfile.ClassFile) {
	jvm.classMutex.Lock()
	defer jvm.classMutex.Unlock()
	jvm.classCache[name] = cf
}

// GetClass retrieves a loaded class
func (jvm *JVM) GetClass(name string) *classfile.ClassFile {
	jvm.classMutex.RLock()
	defer jvm.classMutex.RUnlock()
	return jvm.classCache[name]
}

// GetOrCreateMonitor gets or creates a monitor for an object
func (jvm *JVM) GetOrCreateMonitor(obj any) *Monitor {
	jvm.monitorMutex.Lock()
	defer jvm.monitorMutex.Unlock()

	if monitor, exists := jvm.monitors[obj]; exists {
		return monitor
	}

	monitor := &Monitor{
		waitSet: make([]*Thread, 0),
	}
	monitor.cond = sync.NewCond(&monitor.mutex)
	jvm.monitors[obj] = monitor
	return monitor
}

// Enter acquires the monitor
func (m *Monitor) Enter(thread *Thread) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// If no owner or we already own it (reentrant)
	if m.owner == nil || m.owner == thread {
		m.owner = thread
		m.entryCount++
		return
	}

	// Wait for the monitor to be released
	for m.owner != nil && m.owner != thread {
		m.cond.Wait()
	}

	m.owner = thread
	m.entryCount++
}

// Exit releases the monitor
func (m *Monitor) Exit(thread *Thread) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.owner != thread {
		return fmt.Errorf("IllegalMonitorStateException: not owner of monitor")
	}

	m.entryCount--
	if m.entryCount == 0 {
		m.owner = nil
		m.cond.Signal() // Wake up one waiting thread
	}

	return nil
}

// Wait causes the current thread to wait on this monitor
func (m *Monitor) Wait(thread *Thread) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.owner != thread {
		return fmt.Errorf("IllegalMonitorStateException: not owner of monitor")
	}

	// Save entry count and release monitor
	savedCount := m.entryCount
	m.entryCount = 0
	m.owner = nil
	m.waitSet = append(m.waitSet, thread)
	m.cond.Signal() // Let another thread acquire

	// Wait to be notified
	m.cond.Wait()

	// Reacquire monitor
	for m.owner != nil && m.owner != thread {
		m.cond.Wait()
	}
	m.owner = thread
	m.entryCount = savedCount

	// Remove from wait set
	for i, t := range m.waitSet {
		if t == thread {
			m.waitSet = append(m.waitSet[:i], m.waitSet[i+1:]...)
			break
		}
	}

	return nil
}

// Notify wakes up one waiting thread
func (m *Monitor) Notify(thread *Thread) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.owner != thread {
		return fmt.Errorf("IllegalMonitorStateException: not owner of monitor")
	}

	if len(m.waitSet) > 0 {
		m.cond.Signal()
	}

	return nil
}

// NotifyAll wakes up all waiting threads
func (m *Monitor) NotifyAll(thread *Thread) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.owner != thread {
		return fmt.Errorf("IllegalMonitorStateException: not owner of monitor")
	}

	m.cond.Broadcast()
	return nil
}

// IsRunning returns true if the JVM is still running
func (jvm *JVM) IsRunning() bool {
	return jvm.running.Load()
}

// Shutdown stops the JVM
func (jvm *JVM) Shutdown() {
	jvm.running.Store(false)
}

// GetHeap returns the JVM heap
func (jvm *JVM) GetHeap() *Heap {
	return jvm.heap
}
