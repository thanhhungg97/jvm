package interpreter

import (
	"fmt"
	"simplejvm/classfile"
	"simplejvm/runtime"
)

// loadConstant loads a constant from the constant pool
func (i *Interpreter) loadConstant(frame *runtime.Frame, index uint16) {
	cp := frame.Class.ConstantPool
	c := cp[index]

	switch constant := c.(type) {
	case *classfile.ConstantIntegerInfo:
		frame.OperandStack.PushInt(constant.Value)
	case *classfile.ConstantFloatInfo:
		frame.OperandStack.PushInt(int32(constant.Value)) // Simplified
	case *classfile.ConstantStringInfo:
		str := cp.GetUtf8(constant.StringIndex)
		frame.OperandStack.PushRef(str)
	case *classfile.ConstantClassInfo:
		frame.OperandStack.PushRef(cp.GetUtf8(constant.NameIndex))
	}
}

// loadConstant2 loads a 2-slot constant (long/double)
func (i *Interpreter) loadConstant2(frame *runtime.Frame, index uint16) {
	cp := frame.Class.ConstantPool
	c := cp[index]

	switch constant := c.(type) {
	case *classfile.ConstantLongInfo:
		frame.OperandStack.PushLong(constant.Value)
	case *classfile.ConstantDoubleInfo:
		frame.OperandStack.PushLong(int64(constant.Value)) // Simplified
	}
}

// handlePrintln handles System.out.println calls
func (i *Interpreter) handlePrintln(frame *runtime.Frame, descriptor string) {
	stack := frame.OperandStack

	switch descriptor {
	case "()V":
		stack.PopRef()
		fmt.Println()
	case "(I)V":
		val := stack.PopInt()
		stack.PopRef()
		fmt.Println(val)
	case "(J)V":
		val := stack.PopLong()
		stack.PopRef()
		fmt.Println(val)
	case "(Z)V":
		val := stack.PopInt()
		stack.PopRef()
		if val != 0 {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	case "(C)V":
		val := stack.PopInt()
		stack.PopRef()
		fmt.Println(string(rune(val)))
	case "(Ljava/lang/String;)V":
		val := stack.PopRef()
		stack.PopRef()
		if str, ok := val.(string); ok {
			fmt.Println(str)
		} else {
			fmt.Println(val)
		}
	case "(Ljava/lang/Object;)V":
		val := stack.PopRef()
		stack.PopRef()
		fmt.Println(val)
	default:
		stack.PopSlot()
		stack.PopRef()
		fmt.Println("<unknown println>")
	}
}

// handlePrint handles System.out.print calls (without newline)
func (i *Interpreter) handlePrint(frame *runtime.Frame, descriptor string) {
	stack := frame.OperandStack

	switch descriptor {
	case "(I)V":
		val := stack.PopInt()
		stack.PopRef()
		fmt.Print(val)
	case "(J)V":
		val := stack.PopLong()
		stack.PopRef()
		fmt.Print(val)
	case "(Z)V":
		val := stack.PopInt()
		stack.PopRef()
		if val != 0 {
			fmt.Print("true")
		} else {
			fmt.Print("false")
		}
	case "(C)V":
		val := stack.PopInt()
		stack.PopRef()
		fmt.Print(string(rune(val)))
	case "(Ljava/lang/String;)V":
		val := stack.PopRef()
		stack.PopRef()
		if str, ok := val.(string); ok {
			fmt.Print(str)
		} else {
			fmt.Print(val)
		}
	default:
		stack.PopSlot()
		stack.PopRef()
		fmt.Print("<unknown>")
	}
}

// handleException handles an exception by looking for a handler in the call stack
func (i *Interpreter) handleException(exRef any, exClassName string) error {
	for !i.thread.IsStackEmpty() {
		frame := i.thread.CurrentFrame()

		code := frame.Method.GetCodeAttribute(frame.Class.ConstantPool)
		if code != nil {
			handlerPC := runtime.FindExceptionHandler(code, frame.Class.ConstantPool, frame.PC-1, exClassName)
			if handlerPC >= 0 {
				frame.PC = handlerPC
				frame.OperandStack.Clear()
				frame.OperandStack.PushRef(exRef)
				return nil
			}
		}

		i.thread.PopFrame()
	}

	return fmt.Errorf("uncaught exception: %s", exClassName)
}

// traceCall prints a function call with arguments
func (i *Interpreter) traceCall(methodName string, args []int32) {
	if !i.trace {
		return
	}
	if i.traceMethod != "" && i.traceMethod != methodName {
		return
	}

	depth := i.thread.StackDepth()
	indent := ""
	for j := 0; j < depth-1; j++ {
		indent += "  "
	}

	argsStr := ""
	for idx, arg := range args {
		if idx > 0 {
			argsStr += ", "
		}
		argsStr += fmt.Sprintf("%d", arg)
	}

	fmt.Printf("%s→ %s(%s)\n", indent, methodName, argsStr)
}

// traceReturn prints a function return with value
func (i *Interpreter) traceReturn(methodName string, returnVal int32, hasReturn bool) {
	if !i.trace {
		return
	}
	if i.traceMethod != "" && i.traceMethod != methodName {
		return
	}

	depth := i.thread.StackDepth() + 1
	indent := ""
	for j := 0; j < depth-1; j++ {
		indent += "  "
	}

	if hasReturn {
		fmt.Printf("%s← %s = %d\n", indent, methodName, returnVal)
	} else {
		fmt.Printf("%s← %s\n", indent, methodName)
	}
}

// parseArgTypes returns the type character for each argument in a method descriptor
func parseArgTypes(descriptor string) []byte {
	var types []byte
	i := 1 // Skip '('
	for i < len(descriptor) && descriptor[i] != ')' {
		switch descriptor[i] {
		case 'B', 'C', 'F', 'I', 'S', 'Z':
			types = append(types, descriptor[i])
			i++
		case 'D', 'J':
			types = append(types, descriptor[i])
			i++
		case 'L':
			types = append(types, 'L')
			for descriptor[i] != ';' {
				i++
			}
			i++
		case '[':
			types = append(types, '[')
			i++
			if descriptor[i] == 'L' {
				for descriptor[i] != ';' {
					i++
				}
			}
			i++
		default:
			i++
		}
	}
	return types
}

// countArgs counts the number of argument slots from a method descriptor
func countArgs(descriptor string) int {
	count := 0
	i := 1
	for i < len(descriptor) && descriptor[i] != ')' {
		switch descriptor[i] {
		case 'B', 'C', 'F', 'I', 'S', 'Z':
			count++
			i++
		case 'D', 'J':
			count++
			i++
		case 'L':
			count++
			for descriptor[i] != ';' {
				i++
			}
			i++
		case '[':
			count++
			i++
			if descriptor[i] == 'L' {
				for descriptor[i] != ';' {
					i++
				}
			}
			i++
		default:
			i++
		}
	}
	return count
}

// extractClassName extracts the class name from a placeholder string like "Object<className>"
func extractClassName(s string) string {
	if len(s) > 7 && s[:7] == "Object<" && s[len(s)-1] == '>' {
		return s[7 : len(s)-1]
	}
	return s
}

// arrayTypeChar returns the JVM type character for an array type
func arrayTypeChar(t runtime.ArrayType) rune {
	switch t {
	case runtime.ArrayTypeBoolean:
		return 'Z'
	case runtime.ArrayTypeByte:
		return 'B'
	case runtime.ArrayTypeChar:
		return 'C'
	case runtime.ArrayTypeShort:
		return 'S'
	case runtime.ArrayTypeInt:
		return 'I'
	case runtime.ArrayTypeLong:
		return 'J'
	case runtime.ArrayTypeFloat:
		return 'F'
	case runtime.ArrayTypeDouble:
		return 'D'
	default:
		return '?'
	}
}

// tryLoadClass attempts to lazy-load a class from the same directory as the current class
func (i *Interpreter) tryLoadClass(className string, currentClass *classfile.ClassFile) *classfile.ClassFile {
	// Try to find the class file in the same directory as the current class
	// This is a simple heuristic - look for ClassName.class in examples/

	// Common locations to try
	paths := []string{
		className + ".class",
		"examples/" + className + ".class",
	}

	for _, path := range paths {
		cf, err := classfile.ParseFile(path)
		if err == nil {
			// Successfully loaded - cache it
			i.thread.LoadClass(className, cf)
			return cf
		}
	}

	return nil
}
