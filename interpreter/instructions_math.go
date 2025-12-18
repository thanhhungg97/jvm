package interpreter

import "simplejvm/runtime"

// executeMathInstruction handles arithmetic, bitwise, and conversion instructions
func (i *Interpreter) executeMathInstruction(frame *runtime.Frame, opcode uint8) (bool, error) {
	stack := frame.OperandStack
	locals := frame.LocalVars

	switch opcode {
	// Stack manipulation
	case POP:
		stack.Pop()
	case POP2:
		stack.Pop()
		stack.Pop()
	case DUP:
		stack.Dup()
	case SWAP:
		stack.Swap()

	// Integer arithmetic
	case IADD:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		stack.PushInt(v1 + v2)
	case ISUB:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		stack.PushInt(v1 - v2)
	case IMUL:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		stack.PushInt(v1 * v2)
	case IDIV:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v2 == 0 {
			if err := i.handleException("java/lang/ArithmeticException", "java/lang/ArithmeticException"); err != nil {
				return true, err
			}
			return true, nil
		}
		stack.PushInt(v1 / v2)
	case IREM:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v2 == 0 {
			if err := i.handleException("java/lang/ArithmeticException", "java/lang/ArithmeticException"); err != nil {
				return true, err
			}
			return true, nil
		}
		stack.PushInt(v1 % v2)
	case INEG:
		stack.PushInt(-stack.PopInt())

	// Long arithmetic
	case LADD:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		stack.PushLong(v1 + v2)
	case LSUB:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		stack.PushLong(v1 - v2)
	case LMUL:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		stack.PushLong(v1 * v2)
	case LDIV:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		if v2 == 0 {
			if err := i.handleException("java/lang/ArithmeticException", "java/lang/ArithmeticException"); err != nil {
				return true, err
			}
			return true, nil
		}
		stack.PushLong(v1 / v2)
	case LREM:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		if v2 == 0 {
			if err := i.handleException("java/lang/ArithmeticException", "java/lang/ArithmeticException"); err != nil {
				return true, err
			}
			return true, nil
		}
		stack.PushLong(v1 % v2)
	case LNEG:
		stack.PushLong(-stack.PopLong())

	// Bitwise operations
	case IAND:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		stack.PushInt(v1 & v2)
	case IOR:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		stack.PushInt(v1 | v2)
	case IXOR:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		stack.PushInt(v1 ^ v2)
	case LAND:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		stack.PushLong(v1 & v2)
	case LOR:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		stack.PushLong(v1 | v2)
	case LXOR:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		stack.PushLong(v1 ^ v2)

	// Shifts
	case ISHL:
		v2 := stack.PopInt() & 0x1f
		v1 := stack.PopInt()
		stack.PushInt(v1 << v2)
	case ISHR:
		v2 := stack.PopInt() & 0x1f
		v1 := stack.PopInt()
		stack.PushInt(v1 >> v2)
	case IUSHR:
		v2 := stack.PopInt() & 0x1f
		v1 := stack.PopInt()
		stack.PushInt(int32(uint32(v1) >> v2))

	// Increment local variable
	case IINC:
		index := frame.ReadU1()
		constVal := frame.ReadI1()
		locals.SetInt(int(index), locals.GetInt(int(index))+int32(constVal))

	// Conversions
	case I2L:
		stack.PushLong(int64(stack.PopInt()))
	case L2I:
		stack.PushInt(int32(stack.PopLong()))

	// Long compare
	case LCMP:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		if v1 > v2 {
			stack.PushInt(1)
		} else if v1 < v2 {
			stack.PushInt(-1)
		} else {
			stack.PushInt(0)
		}

	default:
		return false, nil
	}
	return true, nil
}
