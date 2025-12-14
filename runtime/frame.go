package runtime

import "simplejvm/classfile"

// Frame represents a stack frame for method execution
type Frame struct {
	LocalVars    LocalVars
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
		LocalVars:    make(LocalVars, code.MaxLocals),
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

// LocalVars represents local variables
type LocalVars []int64

// SetInt sets an int value
func (l LocalVars) SetInt(index int, val int32) {
	l[index] = int64(val)
}

// GetInt gets an int value
func (l LocalVars) GetInt(index int) int32 {
	return int32(l[index])
}

// SetLong sets a long value
func (l LocalVars) SetLong(index int, val int64) {
	l[index] = val
}

// GetLong gets a long value
func (l LocalVars) GetLong(index int) int64 {
	return l[index]
}

// SetRef sets a reference (stored as int64 for simplicity)
func (l LocalVars) SetRef(index int, val int64) {
	l[index] = val
}

// GetRef gets a reference
func (l LocalVars) GetRef(index int) int64 {
	return l[index]
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
	s.slots[s.size] = int64(val)
	s.size++
}

// PopInt pops an int value
func (s *OperandStack) PopInt() int32 {
	s.size--
	return int32(s.slots[s.size])
}

// PushLong pushes a long value
func (s *OperandStack) PushLong(val int64) {
	s.slots[s.size] = val
	s.size++
}

// PopLong pops a long value
func (s *OperandStack) PopLong() int64 {
	s.size--
	return s.slots[s.size]
}

// PushRef pushes a reference
func (s *OperandStack) PushRef(val interface{}) {
	s.refs[s.size] = val
	s.size++
}

// PopRef pops a reference
func (s *OperandStack) PopRef() interface{} {
	s.size--
	ref := s.refs[s.size]
	s.refs[s.size] = nil
	return ref
}

// PushSlot pushes a raw slot value
func (s *OperandStack) PushSlot(val int64) {
	s.slots[s.size] = val
	s.size++
}

// PopSlot pops a raw slot value
func (s *OperandStack) PopSlot() int64 {
	s.size--
	return s.slots[s.size]
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
