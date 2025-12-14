package interpreter

import (
	"fmt"
	"simplejvm/classfile"
	"simplejvm/runtime"
)

// Interpreter executes bytecode
type Interpreter struct {
	thread      *runtime.Thread
	verbose     bool
	trace       bool   // Simple call trace
	traceMethod string // Only trace this method (empty = all)
}

// NewInterpreter creates a new interpreter
func NewInterpreter(verbose bool) *Interpreter {
	return &Interpreter{
		thread:  runtime.NewThread(),
		verbose: verbose,
		trace:   false,
	}
}

// SetTrace enables simple call/return tracing for a specific method
func (i *Interpreter) SetTrace(methodName string) {
	i.trace = true
	i.traceMethod = methodName
}

// Execute runs the main method of a class
func (i *Interpreter) Execute(cf *classfile.ClassFile) error {
	// Load the class
	className := cf.ClassName()
	i.thread.LoadClass(className, cf)

	// Find the main method
	mainMethod := cf.GetMethod("main", "([Ljava/lang/String;)V")
	if mainMethod == nil {
		return fmt.Errorf("main method not found in class %s", className)
	}

	// Create initial frame
	frame := runtime.NewFrame(i.thread, mainMethod, cf)
	if frame == nil {
		return fmt.Errorf("could not create frame for main method")
	}

	// Push null for args (we don't support command line args yet)
	frame.LocalVars.SetRef(0, 0)

	i.thread.PushFrame(frame)
	return i.run()
}

// ExecuteMethod runs a specific method
func (i *Interpreter) ExecuteMethod(cf *classfile.ClassFile, methodName, descriptor string) error {
	className := cf.ClassName()
	i.thread.LoadClass(className, cf)

	method := cf.GetMethod(methodName, descriptor)
	if method == nil {
		return fmt.Errorf("method %s%s not found", methodName, descriptor)
	}

	frame := runtime.NewFrame(i.thread, method, cf)
	if frame == nil {
		return fmt.Errorf("could not create frame for method %s", methodName)
	}

	i.thread.PushFrame(frame)
	return i.run()
}

// run executes the bytecode loop
func (i *Interpreter) run() error {
	for !i.thread.IsStackEmpty() {
		frame := i.thread.CurrentFrame()
		if frame == nil {
			break
		}

		if frame.PC >= len(frame.Code) {
			// End of method, pop frame
			i.thread.PopFrame()
			continue
		}

		opcode := frame.ReadU1()
		methodName := frame.Method.Name(frame.Class.ConstantPool)

		if i.verbose {
			fmt.Printf("[%s] PC=%d opcode=0x%02X\n", methodName, frame.PC-1, opcode)
		}

		if err := i.executeInstruction(frame, opcode); err != nil {
			return err
		}
	}
	return nil
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

	depth := i.thread.StackDepth() + 1 // +1 because we already popped
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

// executeInstruction executes a single bytecode instruction
func (i *Interpreter) executeInstruction(frame *runtime.Frame, opcode uint8) error {
	stack := frame.OperandStack
	locals := frame.LocalVars
	cp := frame.Class.ConstantPool

	switch opcode {
	// NOP
	case NOP:
		// Do nothing

	// Push null reference
	case ACONST_NULL:
		stack.PushRef(nil)

	// Push int constants
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

	// Push long constants
	case LCONST_0:
		stack.PushLong(0)
	case LCONST_1:
		stack.PushLong(1)

	// Push byte as int
	case BIPUSH:
		val := frame.ReadI1()
		stack.PushInt(int32(val))

	// Push short as int
	case SIPUSH:
		val := frame.ReadI2()
		stack.PushInt(int32(val))

	// Load constant from pool
	case LDC:
		index := frame.ReadU1()
		i.loadConstant(frame, uint16(index))
	case LDC_W:
		index := frame.ReadU2()
		i.loadConstant(frame, index)
	case LDC2_W:
		index := frame.ReadU2()
		i.loadConstant2(frame, index)

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
		stack.PushSlot(locals.GetRef(int(index)))
	case ALOAD_0:
		stack.PushSlot(locals.GetRef(0))
	case ALOAD_1:
		stack.PushSlot(locals.GetRef(1))
	case ALOAD_2:
		stack.PushSlot(locals.GetRef(2))
	case ALOAD_3:
		stack.PushSlot(locals.GetRef(3))

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
		locals.SetRef(int(index), stack.PopSlot())
	case ASTORE_0:
		locals.SetRef(0, stack.PopSlot())
	case ASTORE_1:
		locals.SetRef(1, stack.PopSlot())
	case ASTORE_2:
		locals.SetRef(2, stack.PopSlot())
	case ASTORE_3:
		locals.SetRef(3, stack.PopSlot())

	// Stack manipulation
	case POP:
		stack.PopSlot()
	case POP2:
		stack.PopSlot()
		stack.PopSlot()
	case DUP:
		val := stack.PopSlot()
		stack.PushSlot(val)
		stack.PushSlot(val)
	case SWAP:
		v1 := stack.PopSlot()
		v2 := stack.PopSlot()
		stack.PushSlot(v1)
		stack.PushSlot(v2)

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
			return fmt.Errorf("ArithmeticException: division by zero")
		}
		stack.PushInt(v1 / v2)
	case IREM:
		v2 := stack.PopInt()
		v1 := stack.PopInt()
		if v2 == 0 {
			return fmt.Errorf("ArithmeticException: division by zero")
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
			return fmt.Errorf("ArithmeticException: division by zero")
		}
		stack.PushLong(v1 / v2)
	case LREM:
		v2 := stack.PopLong()
		v1 := stack.PopLong()
		if v2 == 0 {
			return fmt.Errorf("ArithmeticException: division by zero")
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

	// Get static field (simplified - just for System.out)
	case GETSTATIC:
		index := frame.ReadU2()
		fieldRef := cp[index].(*classfile.ConstantFieldrefInfo)
		className := cp.GetClassName(fieldRef.ClassIndex)
		fieldName, _ := cp.GetNameAndType(fieldRef.NameAndTypeIndex)

		if className == "java/lang/System" && fieldName == "out" {
			// Push a placeholder for System.out
			stack.PushRef("System.out")
		} else {
			stack.PushRef(nil)
		}

	// Invoke virtual method
	case INVOKEVIRTUAL:
		index := frame.ReadU2()
		methodRef := cp[index].(*classfile.ConstantMethodrefInfo)
		className := cp.GetClassName(methodRef.ClassIndex)
		methodName, descriptor := cp.GetNameAndType(methodRef.NameAndTypeIndex)

		// Handle System.out.println
		if className == "java/io/PrintStream" && methodName == "println" {
			i.handlePrintln(frame, descriptor)
		} else {
			// For other virtual calls, we'd need full object support
			return fmt.Errorf("unsupported invokevirtual: %s.%s%s", className, methodName, descriptor)
		}

	// Invoke static method
	case INVOKESTATIC:
		index := frame.ReadU2()
		methodRef := cp[index].(*classfile.ConstantMethodrefInfo)
		className := cp.GetClassName(methodRef.ClassIndex)
		methodName, descriptor := cp.GetNameAndType(methodRef.NameAndTypeIndex)

		if err := i.invokeStatic(frame, className, methodName, descriptor); err != nil {
			return err
		}

	// Invoke special (constructors, super calls, private methods)
	case INVOKESPECIAL:
		index := frame.ReadU2()
		methodRef := cp[index].(*classfile.ConstantMethodrefInfo)
		className := cp.GetClassName(methodRef.ClassIndex)
		methodName, descriptor := cp.GetNameAndType(methodRef.NameAndTypeIndex)

		// For now, just pop the object reference for constructors
		if methodName == "<init>" {
			stack.PopRef() // Pop 'this'
			// Pop arguments based on descriptor
			argCount := countArgs(descriptor)
			for j := 0; j < argCount; j++ {
				stack.PopSlot()
			}
		} else {
			return fmt.Errorf("unsupported invokespecial: %s.%s%s", className, methodName, descriptor)
		}

	// New object
	case NEW:
		index := frame.ReadU2()
		className := cp.GetClassName(index)
		// Push a placeholder object
		stack.PushRef(fmt.Sprintf("Object<%s>", className))

	default:
		return fmt.Errorf("unimplemented opcode: 0x%02X at PC=%d", opcode, frame.PC-1)
	}

	return nil
}

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
		stack.PopRef() // Pop PrintStream ref
		fmt.Println()
	case "(I)V":
		val := stack.PopInt()
		stack.PopRef() // Pop PrintStream ref
		fmt.Println(val)
	case "(J)V":
		val := stack.PopLong()
		stack.PopRef() // Pop PrintStream ref
		fmt.Println(val)
	case "(Z)V":
		val := stack.PopInt()
		stack.PopRef() // Pop PrintStream ref
		if val != 0 {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	case "(C)V":
		val := stack.PopInt()
		stack.PopRef() // Pop PrintStream ref
		fmt.Println(string(rune(val)))
	case "(Ljava/lang/String;)V":
		val := stack.PopRef()
		stack.PopRef() // Pop PrintStream ref
		if str, ok := val.(string); ok {
			fmt.Println(str)
		} else {
			fmt.Println(val)
		}
	case "(Ljava/lang/Object;)V":
		val := stack.PopRef()
		stack.PopRef() // Pop PrintStream ref
		fmt.Println(val)
	default:
		// Pop unknown args
		stack.PopSlot()
		stack.PopRef()
		fmt.Println("<unknown println>")
	}
}

// invokeStatic invokes a static method
func (i *Interpreter) invokeStatic(frame *runtime.Frame, className, methodName, descriptor string) error {
	// Look for the class
	cf := i.thread.GetClass(className)
	if cf == nil {
		// Same class invocation
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
	argCount := countArgs(descriptor)
	args := make([]int32, argCount)
	for j := argCount - 1; j >= 0; j-- {
		val := frame.OperandStack.PopSlot()
		newFrame.LocalVars[j] = val
		args[j] = int32(val)
	}

	i.thread.PushFrame(newFrame)

	// Trace the call
	i.traceCall(methodName, args)

	return nil
}

// countArgs counts the number of argument slots from a method descriptor
func countArgs(descriptor string) int {
	count := 0
	i := 1 // Skip '('
	for i < len(descriptor) && descriptor[i] != ')' {
		switch descriptor[i] {
		case 'B', 'C', 'F', 'I', 'S', 'Z':
			count++
			i++
		case 'D', 'J':
			count++ // Long/double could be 2 slots, but we use 1 for simplicity
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
			// Skip array element type
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
