package interpreter

import "simplejvm/runtime"

// executeControlInstruction handles branch and control flow instructions
func (i *Interpreter) executeControlInstruction(frame *runtime.Frame, opcode uint8) bool {
	stack := frame.OperandStack

	switch opcode {
	// Conditional branches (compare with zero)
	case IFEQ:
		offset := frame.ReadI2()
		if stack.PopInt() == 0 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IFNE:
		offset := frame.ReadI2()
		if stack.PopInt() != 0 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IFLT:
		offset := frame.ReadI2()
		if stack.PopInt() < 0 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IFGE:
		offset := frame.ReadI2()
		if stack.PopInt() >= 0 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IFGT:
		offset := frame.ReadI2()
		if stack.PopInt() > 0 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IFLE:
		offset := frame.ReadI2()
		if stack.PopInt() <= 0 {
			frame.PC = frame.PC - 3 + int(offset)
		}

	// Conditional branches (compare two ints)
	case IF_ICMPEQ:
		offset := frame.ReadI2()
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v1 == v2 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IF_ICMPNE:
		offset := frame.ReadI2()
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v1 != v2 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IF_ICMPLT:
		offset := frame.ReadI2()
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v1 < v2 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IF_ICMPGE:
		offset := frame.ReadI2()
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v1 >= v2 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IF_ICMPGT:
		offset := frame.ReadI2()
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v1 > v2 {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IF_ICMPLE:
		offset := frame.ReadI2()
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v1 <= v2 {
			frame.PC = frame.PC - 3 + int(offset)
		}

	// Null checks
	case IFNULL:
		offset := frame.ReadI2()
		ref := stack.PopRef()
		if ref == nil {
			frame.PC = frame.PC - 3 + int(offset)
		}
	case IFNONNULL:
		offset := frame.ReadI2()
		ref := stack.PopRef()
		if ref != nil {
			frame.PC = frame.PC - 3 + int(offset)
		}

	// Unconditional jump
	case GOTO:
		offset := frame.ReadI2()
		frame.PC = frame.PC - 3 + int(offset)
	case GOTO_W:
		offset := frame.ReadI4()
		frame.PC = frame.PC - 5 + int(offset)

	// Return from method
	case RETURN:
		i.thread.PopFrame()
	case IRETURN:
		retVal := stack.PopInt()
		methodName := frame.Method.Name(frame.Class.ConstantPool)
		i.traceReturn(methodName, retVal, true)
		i.thread.PopFrame()
		caller := i.thread.CurrentFrame()
		if caller != nil {
			caller.OperandStack.PushInt(retVal)
		}
	case LRETURN:
		retVal := stack.PopLong()
		i.thread.PopFrame()
		caller := i.thread.CurrentFrame()
		if caller != nil {
			caller.OperandStack.PushLong(retVal)
		}
	case ARETURN:
		retVal := stack.PopRef()
		i.thread.PopFrame()
		caller := i.thread.CurrentFrame()
		if caller != nil {
			caller.OperandStack.PushRef(retVal)
		}

	default:
		return false
	}
	return true
}
