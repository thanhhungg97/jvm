package interpreter

import (
	"fmt"
	"simplejvm/classfile"
	"simplejvm/runtime"
)

// executeInvokeInstruction handles method invocation instructions
func (i *Interpreter) executeInvokeInstruction(frame *runtime.Frame, opcode uint8) (bool, error) {
	stack := frame.OperandStack
	cp := frame.Class.ConstantPool

	switch opcode {
	// Invoke virtual method
	case INVOKEVIRTUAL:
		index := frame.ReadU2()
		methodRef := cp[index].(*classfile.ConstantMethodrefInfo)
		className := cp.GetClassName(methodRef.ClassIndex)
		methodName, descriptor := cp.GetNameAndType(methodRef.NameAndTypeIndex)

		// Handle System.out.println
		if className == "java/io/PrintStream" && methodName == "println" {
			i.handlePrintln(frame, descriptor)
		} else if className == "java/io/PrintStream" && methodName == "print" {
			i.handlePrint(frame, descriptor)
		} else {
			if err := i.invokeVirtual(frame, className, methodName, descriptor); err != nil {
				return true, err
			}
		}

	// Invoke static method
	case INVOKESTATIC:
		index := frame.ReadU2()
		methodRef := cp[index].(*classfile.ConstantMethodrefInfo)
		className := cp.GetClassName(methodRef.ClassIndex)
		methodName, descriptor := cp.GetNameAndType(methodRef.NameAndTypeIndex)

		if err := i.invokeStatic(frame, className, methodName, descriptor); err != nil {
			return true, err
		}

	// Invoke special (constructors, super calls, private methods)
	case INVOKESPECIAL:
		index := frame.ReadU2()
		methodRef := cp[index].(*classfile.ConstantMethodrefInfo)
		className := cp.GetClassName(methodRef.ClassIndex)
		methodName, descriptor := cp.GetNameAndType(methodRef.NameAndTypeIndex)

		targetClass := i.thread.GetClass(className)
		if targetClass == nil && className == frame.Class.ClassName() {
			targetClass = frame.Class
		}
		// Try lazy class loading
		if targetClass == nil {
			targetClass = i.tryLoadClass(className, frame.Class)
		}

		if targetClass != nil {
			method := targetClass.GetMethod(methodName, descriptor)
			if method != nil {
				code := method.GetCodeAttribute(targetClass.ConstantPool)
				if code != nil {
					newFrame := runtime.NewFrame(i.thread, method, targetClass)
					argTypes := parseArgTypes(descriptor)
					argCount := len(argTypes)

					// Pop arguments in reverse order, correctly handling types
					for j := argCount - 1; j >= 0; j-- {
						if argTypes[j] == 'L' || argTypes[j] == '[' {
							newFrame.LocalVars.SetRef(j+1, stack.PopRef())
						} else {
							newFrame.LocalVars.SetSlot(j+1, stack.PopSlot())
						}
					}
					// Pop 'this' reference
					newFrame.LocalVars.SetRef(0, stack.PopRef())

					i.thread.PushFrame(newFrame)
					return true, nil
				}
			}
		}

		// Fallback: just pop arguments for constructors we can't execute
		if methodName == "<init>" {
			argCount := countArgs(descriptor)
			for j := 0; j < argCount; j++ {
				stack.PopSlot()
			}
			stack.PopRef() // Pop 'this'
		} else {
			return true, fmt.Errorf("unsupported invokespecial: %s.%s%s", className, methodName, descriptor)
		}

	default:
		return false, nil
	}
	return true, nil
}

// invokeStatic invokes a static method
func (i *Interpreter) invokeStatic(frame *runtime.Frame, className, methodName, descriptor string) error {
	// Check for native method first
	if nativeMethod := runtime.Natives.Lookup(className, methodName, descriptor); nativeMethod != nil {
		return nativeMethod(frame)
	}

	// Look for the class
	cf := i.thread.GetClass(className)
	if cf == nil {
		if className == frame.Class.ClassName() {
			cf = frame.Class
		} else {
			return fmt.Errorf("class not found: %s", className)
		}
	}

	method := cf.GetMethod(methodName, descriptor)
	if method == nil {
		return fmt.Errorf("method not found: %s.%s%s", className, methodName, descriptor)
	}

	newFrame := runtime.NewFrame(i.thread, method, cf)
	if newFrame == nil {
		return fmt.Errorf("could not create frame for %s.%s", className, methodName)
	}

	// Pass arguments from caller's stack to new frame's locals
	argTypes := parseArgTypes(descriptor)
	argCount := len(argTypes)
	args := make([]int32, argCount)
	for j := argCount - 1; j >= 0; j-- {
		if argTypes[j] == 'L' || argTypes[j] == '[' {
			ref := frame.OperandStack.PopRef()
			newFrame.LocalVars.SetRef(j, ref)
			args[j] = 0
		} else {
			val := frame.OperandStack.PopSlot()
			newFrame.LocalVars.SetSlot(j, val)
			args[j] = int32(val)
		}
	}

	i.thread.PushFrame(newFrame)
	i.traceCall(methodName, args)

	return nil
}

// invokeVirtual invokes an instance method
func (i *Interpreter) invokeVirtual(frame *runtime.Frame, className, methodName, descriptor string) error {
	stack := frame.OperandStack

	argTypes := parseArgTypes(descriptor)
	argCount := len(argTypes)

	type argValue struct {
		slot int64
		ref  any
	}
	args := make([]argValue, argCount)
	for j := argCount - 1; j >= 0; j-- {
		if argTypes[j] == 'L' || argTypes[j] == '[' {
			args[j] = argValue{ref: stack.PopRef()}
		} else {
			args[j] = argValue{slot: stack.PopSlot()}
		}
	}
	objRef := stack.PopRef()

	if objRef == nil {
		return fmt.Errorf("NullPointerException: invokevirtual on null object")
	}

	var targetClass *classfile.ClassFile
	if obj, ok := objRef.(*runtime.Object); ok {
		targetClass = obj.Class
	} else {
		return fmt.Errorf("invokevirtual: not an Object: %T", objRef)
	}

	if targetClass == nil {
		return fmt.Errorf("invokevirtual: object has no class")
	}

	method := targetClass.GetMethod(methodName, descriptor)
	if method == nil {
		return fmt.Errorf("method not found: %s.%s%s", targetClass.ClassName(), methodName, descriptor)
	}

	newFrame := runtime.NewFrame(i.thread, method, targetClass)
	if newFrame == nil {
		return fmt.Errorf("could not create frame for %s.%s", targetClass.ClassName(), methodName)
	}

	newFrame.LocalVars.SetRef(0, objRef)

	for j := 0; j < argCount; j++ {
		if argTypes[j] == 'L' || argTypes[j] == '[' {
			newFrame.LocalVars.SetRef(j+1, args[j].ref)
		} else {
			newFrame.LocalVars.SetSlot(j+1, args[j].slot)
		}
	}

	i.thread.PushFrame(newFrame)
	i.traceCall(methodName, nil)

	return nil
}
