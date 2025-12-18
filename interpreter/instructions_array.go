package interpreter

import (
	"fmt"
	"simplejvm/runtime"
)

// executeArrayInstruction handles array-related instructions
func (i *Interpreter) executeArrayInstruction(frame *runtime.Frame, opcode uint8) (bool, error) {
	stack := frame.OperandStack
	cp := frame.Class.ConstantPool

	switch opcode {
	// Array load instructions
	case IALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushInt(arr.GetInt(index))

	case LALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushLong(arr.GetLong(index))

	case AALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushRef(arr.GetRef(index))

	case BALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushInt(arr.GetInt(index))

	case CALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushInt(arr.GetInt(index))

	case SALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushInt(arr.GetInt(index))

	case FALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushFloat(arr.GetFloat(index))

	case DALOAD:
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		stack.PushDouble(arr.GetDouble(index))

	// Array store instructions
	case IASTORE:
		val := stack.PopInt()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetInt(index, val)

	case LASTORE:
		val := stack.PopLong()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetLong(index, val)

	case AASTORE:
		val := stack.PopRef()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetRef(index, val)

	case BASTORE:
		val := stack.PopInt()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetInt(index, int32(int8(val)))

	case CASTORE:
		val := stack.PopInt()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetInt(index, int32(uint16(val)))

	case SASTORE:
		val := stack.PopInt()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetInt(index, int32(int16(val)))

	case FASTORE:
		val := stack.PopFloat()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetFloat(index, val)

	case DASTORE:
		val := stack.PopDouble()
		index := stack.PopInt()
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		if index < 0 || index >= arr.Length {
			return true, fmt.Errorf("ArrayIndexOutOfBoundsException: %d", index)
		}
		arr.SetDouble(index, val)

	// Create new primitive array
	case NEWARRAY:
		atype := frame.ReadU1()
		count := stack.PopInt()
		if count < 0 {
			return true, fmt.Errorf("NegativeArraySizeException: %d", count)
		}
		arr := runtime.NewPrimitiveArray(runtime.ArrayType(atype), count)
		stack.PushRef(arr)

	// Create new reference array
	case ANEWARRAY:
		index := frame.ReadU2()
		className := cp.GetClassName(index)
		count := stack.PopInt()
		if count < 0 {
			return true, fmt.Errorf("NegativeArraySizeException: %d", count)
		}
		arr := runtime.NewReferenceArray(className, count)
		stack.PushRef(arr)

	// Get array length
	case ARRAYLENGTH:
		arrRef := stack.PopRef()
		if arrRef == nil {
			return true, fmt.Errorf("NullPointerException: array is null")
		}
		arr := arrRef.(*runtime.Array)
		stack.PushInt(arr.Length)

	default:
		return false, nil
	}
	return true, nil
}
