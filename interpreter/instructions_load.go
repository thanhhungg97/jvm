package interpreter

import "simplejvm/runtime"

// executeLoadInstruction handles load instructions (from local vars to stack)
func (i *Interpreter) executeLoadInstruction(frame *runtime.Frame, opcode uint8) bool {
	stack := frame.OperandStack
	locals := frame.LocalVars

	switch opcode {
	// Load int from local variable
	case ILOAD:
		index := frame.ReadU1()
		stack.PushInt(locals.GetInt(int(index)))
	case ILOAD_0:
		stack.PushInt(locals.GetInt(0))
	case ILOAD_1:
		stack.PushInt(locals.GetInt(1))
	case ILOAD_2:
		stack.PushInt(locals.GetInt(2))
	case ILOAD_3:
		stack.PushInt(locals.GetInt(3))

	// Load long from local variable
	case LLOAD:
		index := frame.ReadU1()
		stack.PushLong(locals.GetLong(int(index)))
	case LLOAD_0:
		stack.PushLong(locals.GetLong(0))
	case LLOAD_1:
		stack.PushLong(locals.GetLong(1))
	case LLOAD_2:
		stack.PushLong(locals.GetLong(2))
	case LLOAD_3:
		stack.PushLong(locals.GetLong(3))

	// Load reference from local variable
	case ALOAD:
		index := frame.ReadU1()
		stack.PushRef(locals.GetRef(int(index)))
	case ALOAD_0:
		stack.PushRef(locals.GetRef(0))
	case ALOAD_1:
		stack.PushRef(locals.GetRef(1))
	case ALOAD_2:
		stack.PushRef(locals.GetRef(2))
	case ALOAD_3:
		stack.PushRef(locals.GetRef(3))

	default:
		return false
	}
	return true
}
