package runtime

import (
	"simplejvm/classfile"
	"unsafe"
)

// Frame represents a stack frame for method execution
type Frame struct {
	LocalVars    *LocalVars
	OperandStack *OperandStack
	Thread       *Thread
	Method       *classfile.MethodInfo
	Class        *classfile.ClassFile
	PC           int // Program counter
	Code         []byte
}

// NewFrame creates a new stack frame
func NewFrame(thread *Thread, method *classfile.MethodInfo, class *classfile.ClassFile) *Frame {
	code := method.GetCodeAttribute(class.ConstantPool)
	if code == nil {
		return nil
	}

	return &Frame{
		LocalVars:    NewLocalVars(code.MaxLocals),
		OperandStack: NewOperandStack(int(code.MaxStack)),
		Thread:       thread,
		Method:       method,
		Class:        class,
		PC:           0,
		Code:         code.Code,
	}
}

// NextPC returns the next program counter value
func (f *Frame) NextPC() int {
	return f.PC
}

// SetNextPC sets the next program counter value
func (f *Frame) SetNextPC(pc int) {
	f.PC = pc
}

// ReadCode reads a byte from the bytecode
func (f *Frame) ReadU1() uint8 {
	code := f.Code[f.PC]
	f.PC++
	return code
}

// ReadCode reads a signed byte from the bytecode
func (f *Frame) ReadI1() int8 {
	return int8(f.ReadU1())
}

// ReadU2 reads a 2-byte unsigned value from bytecode
func (f *Frame) ReadU2() uint16 {
	high := uint16(f.ReadU1())
	low := uint16(f.ReadU1())
	return high<<8 | low
}

// ReadI2 reads a 2-byte signed value from bytecode
func (f *Frame) ReadI2() int16 {
	return int16(f.ReadU2())
}

// ReadI4 reads a 4-byte signed value from bytecode
func (f *Frame) ReadI4() int32 {
	b1 := int32(f.ReadU1())
	b2 := int32(f.ReadU1())
	b3 := int32(f.ReadU1())
	b4 := int32(f.ReadU1())
	return b1<<24 | b2<<16 | b3<<8 | b4
}

// LocalVars represents local variables with support for both primitives and references
type LocalVars struct {
	slots []int64
	refs  []any
}

// NewLocalVars creates a new local variables array
func NewLocalVars(maxLocals uint16) *LocalVars {
	return &LocalVars{
		slots: make([]int64, maxLocals),
		refs:  make([]any, maxLocals),
	}
}

// SetInt sets an int value
func (l *LocalVars) SetInt(index int, val int32) {
	l.slots[index] = int64(val)
	l.refs[index] = nil // Clear any reference
}

// GetInt gets an int value
func (l *LocalVars) GetInt(index int) int32 {
	return int32(l.slots[index])
}

// SetLong sets a long value
func (l *LocalVars) SetLong(index int, val int64) {
	l.slots[index] = val
	l.refs[index] = nil
}

// GetLong gets a long value
func (l *LocalVars) GetLong(index int) int64 {
	return l.slots[index]
}

// SetSlot sets a raw slot value (for backwards compatibility)
func (l *LocalVars) SetSlot(index int, val int64) {
	l.slots[index] = val
}

// GetSlot gets a raw slot value
func (l *LocalVars) GetSlot(index int) int64 {
	return l.slots[index]
}

// SetRef sets a reference
func (l *LocalVars) SetRef(index int, val any) {
	l.refs[index] = val
	l.slots[index] = 0
}

// GetRef gets a reference
func (l *LocalVars) GetRef(index int) any {
	return l.refs[index]
}

// OperandStack represents the operand stack
type OperandStack struct {
	size  int
	slots []int64
	refs  []interface{} // For object references
}

// NewOperandStack creates a new operand stack
func NewOperandStack(maxSize int) *OperandStack {
	if maxSize < 1 {
		maxSize = 1
	}
	return &OperandStack{
		size:  0,
		slots: make([]int64, maxSize),
		refs:  make([]interface{}, maxSize),
	}
}

// PushInt pushes an int value
func (s *OperandStack) PushInt(val int32) {
	s.refs[s.size] = nil // Clear ref at this position
	s.slots[s.size] = int64(val)
	s.size++
}

// PopInt pops an int value
func (s *OperandStack) PopInt() int32 {
	s.size--
	s.refs[s.size] = nil
	return int32(s.slots[s.size])
}

// PushLong pushes a long value
func (s *OperandStack) PushLong(val int64) {
	s.refs[s.size] = nil
	s.slots[s.size] = val
	s.size++
}

// PopLong pops a long value
func (s *OperandStack) PopLong() int64 {
	s.size--
	s.refs[s.size] = nil
	return s.slots[s.size]
}

// PushFloat pushes a float value
func (s *OperandStack) PushFloat(val float32) {
	bits := *(*int32)(unsafe.Pointer(&val))
	s.slots[s.size] = int64(bits)
	s.size++
}

// PopFloat pops a float value
func (s *OperandStack) PopFloat() float32 {
	s.size--
	bits := int32(s.slots[s.size])
	return *(*float32)(unsafe.Pointer(&bits))
}

// PushDouble pushes a double value
func (s *OperandStack) PushDouble(val float64) {
	bits := *(*int64)(unsafe.Pointer(&val))
	s.slots[s.size] = bits
	s.size++
}

// PopDouble pops a double value
func (s *OperandStack) PopDouble() float64 {
	s.size--
	bits := s.slots[s.size]
	return *(*float64)(unsafe.Pointer(&bits))
}

// PushRef pushes a reference
func (s *OperandStack) PushRef(val interface{}) {
	s.slots[s.size] = 0 // Clear slot at this position
	s.refs[s.size] = val
	s.size++
}

// PopRef pops a reference
func (s *OperandStack) PopRef() interface{} {
	s.size--
	ref := s.refs[s.size]
	s.refs[s.size] = nil
	s.slots[s.size] = 0
	return ref
}

// PushSlot pushes a raw slot value
func (s *OperandStack) PushSlot(val int64) {
	s.slots[s.size] = val
	s.refs[s.size] = nil // Clear ref at this position
	s.size++
}

// PopSlot pops a raw slot value
func (s *OperandStack) PopSlot() int64 {
	s.size--
	s.refs[s.size] = nil // Clear ref
	return s.slots[s.size]
}

// Pop pops and discards the top value (either slot or ref)
func (s *OperandStack) Pop() {
	s.size--
	s.refs[s.size] = nil
}

// Dup duplicates the top value (handles both slots and refs)
func (s *OperandStack) Dup() {
	ref := s.refs[s.size-1]
	slot := s.slots[s.size-1]
	s.slots[s.size] = slot
	s.refs[s.size] = ref
	s.size++
}

// Swap swaps the top two values
func (s *OperandStack) Swap() {
	s.slots[s.size-1], s.slots[s.size-2] = s.slots[s.size-2], s.slots[s.size-1]
	s.refs[s.size-1], s.refs[s.size-2] = s.refs[s.size-2], s.refs[s.size-1]
}

// Top returns the top value without popping
func (s *OperandStack) TopInt() int32 {
	return int32(s.slots[s.size-1])
}

// Size returns current stack size
func (s *OperandStack) Size() int {
	return s.size
}

// IsEmpty returns true if stack is empty
func (s *OperandStack) IsEmpty() bool {
	return s.size == 0
}

// Clear clears the stack
func (s *OperandStack) Clear() {
	s.size = 0
}

// PeekSlot peeks at a slot value at the given index from bottom (0 = bottom)
func (s *OperandStack) PeekSlot(index int) int64 {
	if index < 0 || index >= s.size {
		return 0
	}
	return s.slots[index]
}

// PeekRef peeks at a reference at the given index from bottom (0 = bottom)
func (s *OperandStack) PeekRef(index int) interface{} {
	if index < 0 || index >= s.size {
		return nil
	}
	return s.refs[index]
}

// HasRefAt returns true if the given index has a reference (vs primitive)
func (s *OperandStack) HasRefAt(index int) bool {
	if index < 0 || index >= s.size {
		return false
	}
	return s.refs[index] != nil
}
