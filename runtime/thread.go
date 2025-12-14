package runtime

import "simplejvm/classfile"

// Thread represents a JVM thread with its stack of frames
type Thread struct {
	stack   []*Frame
	Classes map[string]*classfile.ClassFile // Loaded classes
}

// NewThread creates a new thread
func NewThread() *Thread {
	return &Thread{
		stack:   make([]*Frame, 0, 32),
		Classes: make(map[string]*classfile.ClassFile),
	}
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

// LoadClass loads a class file and caches it
func (t *Thread) LoadClass(name string, cf *classfile.ClassFile) {
	t.Classes[name] = cf
}

// GetClass gets a loaded class
func (t *Thread) GetClass(name string) *classfile.ClassFile {
	return t.Classes[name]
}
