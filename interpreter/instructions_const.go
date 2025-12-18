package interpreter

import "simplejvm/runtime"

// executeConstInstruction handles constant-pushing instructions
func (i *Interpreter) executeConstInstruction(frame *runtime.Frame, opcode uint8) bool {
	stack := frame.OperandStack

	switch opcode {
	case NOP:
		// Do nothing

	case ACONST_NULL:
		stack.PushRef(nil)

	case ICONST_M1:
		stack.PushInt(-1)
	case ICONST_0:
		stack.PushInt(0)
	case ICONST_1:
		stack.PushInt(1)
	case ICONST_2:
		stack.PushInt(2)
	case ICONST_3:
		stack.PushInt(3)
	case ICONST_4:
		stack.PushInt(4)
	case ICONST_5:
		stack.PushInt(5)

	case LCONST_0:
		stack.PushLong(0)
	case LCONST_1:
		stack.PushLong(1)

	case BIPUSH:
		val := frame.ReadI1()
		stack.PushInt(int32(val))

	case SIPUSH:
		val := frame.ReadI2()
		stack.PushInt(int32(val))

	case LDC:
		index := frame.ReadU1()
		i.loadConstant(frame, uint16(index))
	case LDC_W:
		index := frame.ReadU2()
		i.loadConstant(frame, index)
	case LDC2_W:
		index := frame.ReadU2()
		i.loadConstant2(frame, index)

	default:
		return false // Not handled
	}
	return true
}
