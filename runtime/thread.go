package runtime

import "simplejvm/classfile"

// Thread represents a JVM thread with its stack of frames
type Thread struct {
	id      int64                           // Thread ID
	stack   []*Frame                        // Call stack
	Classes map[string]*classfile.ClassFile // Loaded classes
	jvm     *JVM                            // Reference to parent JVM
}

// NewThread creates a new thread (standalone, without JVM)
func NewThread() *Thread {
	return &Thread{
		id:      1,
		stack:   make([]*Frame, 0, 32),
		Classes: make(map[string]*classfile.ClassFile),
	}
}

// ID returns the thread ID
func (t *Thread) ID() int64 {
	return t.id
}

// JVM returns the parent JVM (may be nil)
func (t *Thread) JVM() *JVM {
	return t.jvm
}

// PushFrame pushes a new frame onto the stack
func (t *Thread) PushFrame(frame *Frame) {
	t.stack = append(t.stack, frame)
}

// PopFrame pops a frame from the stack
func (t *Thread) PopFrame() *Frame {
	if len(t.stack) == 0 {
		return nil
	}
	frame := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	return frame
}

// CurrentFrame returns the current frame
func (t *Thread) CurrentFrame() *Frame {
	if len(t.stack) == 0 {
		return nil
	}
	return t.stack[len(t.stack)-1]
}

// IsStackEmpty returns true if the stack is empty
func (t *Thread) IsStackEmpty() bool {
	return len(t.stack) == 0
}

// StackDepth returns the current stack depth
func (t *Thread) StackDepth() int {
	return len(t.stack)
}

// ContainsFrame returns true if the given frame is in the stack
func (t *Thread) ContainsFrame(frame *Frame) bool {
	for _, f := range t.stack {
		if f == frame {
			return true
		}
	}
	return false
}

// LoadClass loads a class file and caches it
func (t *Thread) LoadClass(name string, cf *classfile.ClassFile) {
	t.Classes[name] = cf
}

// GetClass gets a loaded class
func (t *Thread) GetClass(name string) *classfile.ClassFile {
	return t.Classes[name]
}
