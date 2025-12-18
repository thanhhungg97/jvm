package interpreter

import "simplejvm/runtime"

// executeStoreInstruction handles store instructions (from stack to local vars)
func (i *Interpreter) executeStoreInstruction(frame *runtime.Frame, opcode uint8) bool {
	stack := frame.OperandStack
	locals := frame.LocalVars

	switch opcode {
	// Store int to local variable
	case ISTORE:
		index := frame.ReadU1()
		locals.SetInt(int(index), stack.PopInt())
	case ISTORE_0:
		locals.SetInt(0, stack.PopInt())
	case ISTORE_1:
		locals.SetInt(1, stack.PopInt())
	case ISTORE_2:
		locals.SetInt(2, stack.PopInt())
	case ISTORE_3:
		locals.SetInt(3, stack.PopInt())

	// Store long to local variable
	case LSTORE:
		index := frame.ReadU1()
		locals.SetLong(int(index), stack.PopLong())
	case LSTORE_0:
		locals.SetLong(0, stack.PopLong())
	case LSTORE_1:
		locals.SetLong(1, stack.PopLong())
	case LSTORE_2:
		locals.SetLong(2, stack.PopLong())
	case LSTORE_3:
		locals.SetLong(3, stack.PopLong())

	// Store reference to local variable
	case ASTORE:
		index := frame.ReadU1()
		locals.SetRef(int(index), stack.PopRef())
	case ASTORE_0:
		locals.SetRef(0, stack.PopRef())
	case ASTORE_1:
		locals.SetRef(1, stack.PopRef())
	case ASTORE_2:
		locals.SetRef(2, stack.PopRef())
	case ASTORE_3:
		locals.SetRef(3, stack.PopRef())

	default:
		return false
	}
	return true
}
