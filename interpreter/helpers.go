package interpreter

import (
	"fmt"
	"simplejvm/classfile"
	"simplejvm/runtime"
	"strings"
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

// printFrameDebug prints detailed frame state for debugging (simplified format)
func (i *Interpreter) printFrameDebug(frame *runtime.Frame, pc int, opcode uint8) {
	methodName := frame.Method.Name(frame.Class.ConstantPool)
	className := frame.Class.ClassName()

	// Header with method name
	header := fmt.Sprintf("─ %s.%s ", className, methodName)
	fmt.Printf("┌%s%s\n", header, strings.Repeat("─", 60-len(header)))

	// Instruction with description
	opName := getOpcodeName(opcode)
	desc := getOpcodeDescription(frame, pc, opcode)
	fmt.Printf("│ PC=%-3d  %-14s  → %s\n", pc, opName, desc)
	fmt.Println("│")

	// Locals in compact format
	fmt.Printf("│ Locals: %s\n", formatLocalsCompact(frame))

	// Stack in compact format
	fmt.Printf("│ Stack:  %s\n", formatStackCompact(frame))

	fmt.Printf("└%s\n", strings.Repeat("─", 60))
	fmt.Println()
}

// formatLocalsCompact formats local variables in a compact single line
func formatLocalsCompact(frame *runtime.Frame) string {
	code := frame.Method.GetCodeAttribute(frame.Class.ConstantPool)
	maxLocals := 0
	if code != nil {
		maxLocals = int(code.MaxLocals)
	}
	if maxLocals == 0 {
		return "(none)"
	}

	var parts []string
	for j := 0; j < maxLocals && j < 8; j++ { // Limit to 8 for readability
		ref := frame.LocalVars.GetRef(j)
		slot := frame.LocalVars.GetSlot(j)
		if ref != nil {
			parts = append(parts, fmt.Sprintf("[%d]=%s", j, formatRefShort(ref)))
		} else {
			parts = append(parts, fmt.Sprintf("[%d]=%d", j, slot))
		}
	}
	if maxLocals > 8 {
		parts = append(parts, "...")
	}
	return strings.Join(parts, ", ")
}

// formatStackCompact formats operand stack in a compact single line
func formatStackCompact(frame *runtime.Frame) string {
	stackSize := frame.OperandStack.Size()
	if stackSize == 0 {
		return "[]"
	}

	var parts []string
	for j := 0; j < stackSize; j++ {
		if frame.OperandStack.HasRefAt(j) {
			ref := frame.OperandStack.PeekRef(j)
			parts = append(parts, formatRefShort(ref))
		} else {
			slot := frame.OperandStack.PeekSlot(j)
			parts = append(parts, fmt.Sprintf("%d", slot))
		}
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// formatRefShort formats a reference in short form
func formatRefShort(ref any) string {
	if ref == nil {
		return "null"
	}
	switch v := ref.(type) {
	case *runtime.Object:
		// Get short class name
		name := v.ClassName()
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		return fmt.Sprintf("<%s>", name)
	case *runtime.Array:
		return fmt.Sprintf("arr[%d]", v.Length)
	case string:
		if len(v) > 15 {
			return fmt.Sprintf("\"%s...\"", v[:12])
		}
		return fmt.Sprintf("\"%s\"", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// getOpcodeDescription returns a human-readable description of what the instruction does
func getOpcodeDescription(frame *runtime.Frame, pc int, opcode uint8) string {
	// Save current PC to read operands, then we can describe
	code := frame.Code

	switch opcode {
	// Constants
	case 0x00:
		return "Do nothing"
	case 0x01:
		return "Push null"
	case 0x02:
		return "Push -1"
	case 0x03, 0x04, 0x05, 0x06, 0x07, 0x08:
		return fmt.Sprintf("Push %d", int(opcode)-0x03)
	case 0x09:
		return "Push 0L"
	case 0x0A:
		return "Push 1L"
	case 0x10: // bipush
		if pc+1 < len(code) {
			return fmt.Sprintf("Push %d", int8(code[pc+1]))
		}
		return "Push byte"
	case 0x11: // sipush
		if pc+2 < len(code) {
			val := int16(code[pc+1])<<8 | int16(code[pc+2])
			return fmt.Sprintf("Push %d", val)
		}
		return "Push short"
	case 0x12: // ldc
		if pc+1 < len(code) {
			index := uint16(code[pc+1])
			return fmt.Sprintf("Load constant #%d: %s", index, getConstantPoolValue(frame, index))
		}
		return "Load constant from pool"
	case 0x13: // ldc_w
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Load constant #%d: %s", index, getConstantPoolValue(frame, index))
		}
		return "Load constant from pool"

	// Loads
	case 0x15:
		if pc+1 < len(code) {
			return fmt.Sprintf("Load int from local[%d]", code[pc+1])
		}
		return "Load int"
	case 0x1A, 0x1B, 0x1C, 0x1D:
		return fmt.Sprintf("Load int from local[%d]", opcode-0x1A)
	case 0x1E, 0x1F, 0x20, 0x21:
		return fmt.Sprintf("Load long from local[%d]", opcode-0x1E)
	case 0x2A, 0x2B, 0x2C, 0x2D:
		return fmt.Sprintf("Load ref from local[%d]", opcode-0x2A)

	// Stores
	case 0x36:
		if pc+1 < len(code) {
			return fmt.Sprintf("Store int to local[%d]", code[pc+1])
		}
		return "Store int"
	case 0x3B, 0x3C, 0x3D, 0x3E:
		return fmt.Sprintf("Store int to local[%d]", opcode-0x3B)
	case 0x4B, 0x4C, 0x4D, 0x4E:
		return fmt.Sprintf("Store ref to local[%d]", opcode-0x4B)

	// Math
	case 0x60:
		return "Add top two ints"
	case 0x64:
		return "Subtract top two ints"
	case 0x68:
		return "Multiply top two ints"
	case 0x6C:
		return "Divide top two ints"
	case 0x70:
		return "Remainder of top two ints"
	case 0x84: // iinc
		if pc+2 < len(code) {
			return fmt.Sprintf("Increment local[%d] by %d", code[pc+1], int8(code[pc+2]))
		}
		return "Increment local"

	// Stack
	case 0x57:
		return "Pop top value"
	case 0x59:
		return "Duplicate top value"
	case 0x5F:
		return "Swap top two values"

	// Branches
	case 0x99:
		return "Jump if == 0"
	case 0x9A:
		return "Jump if != 0"
	case 0x9B:
		return "Jump if < 0"
	case 0x9C:
		return "Jump if >= 0"
	case 0x9D:
		return "Jump if > 0"
	case 0x9E:
		return "Jump if <= 0"
	case 0x9F:
		return "Jump if equal"
	case 0xA0:
		return "Jump if not equal"
	case 0xA1:
		return "Jump if less than"
	case 0xA2:
		return "Jump if greater or equal"
	case 0xA3:
		return "Jump if greater than"
	case 0xA4:
		return "Jump if less or equal"
	case 0xA7:
		return "Unconditional jump"

	// Returns
	case 0xAC:
		return "Return int"
	case 0xAD:
		return "Return long"
	case 0xB0:
		return "Return reference"
	case 0xB1:
		return "Return void"

	// Fields & Objects
	case 0xB2: // getstatic
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Get static: %s", getConstantPoolValue(frame, index))
		}
		return "Get static field"
	case 0xB3: // putstatic
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Put static: %s", getConstantPoolValue(frame, index))
		}
		return "Put static field"
	case 0xB4: // getfield
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Get field: %s", getConstantPoolValue(frame, index))
		}
		return "Get instance field"
	case 0xB5: // putfield
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Put field: %s", getConstantPoolValue(frame, index))
		}
		return "Put instance field"
	case 0xBB: // new
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			className := frame.Class.ConstantPool.GetClassName(index)
			return fmt.Sprintf("Create new %s", className)
		}
		return "Create new object"

	// Invoke
	case 0xB6: // invokevirtual
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Call: %s", getConstantPoolValue(frame, index))
		}
		return "Call virtual method"
	case 0xB7: // invokespecial
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Call special: %s", getConstantPoolValue(frame, index))
		}
		return "Call special method"
	case 0xB8: // invokestatic
		if pc+2 < len(code) {
			index := uint16(code[pc+1])<<8 | uint16(code[pc+2])
			return fmt.Sprintf("Call: %s", getConstantPoolValue(frame, index))
		}
		return "Call static method"

	// Arrays
	case 0x2E:
		return "Load int from array"
	case 0x32:
		return "Load ref from array"
	case 0x4F:
		return "Store int to array"
	case 0x53:
		return "Store ref to array"
	case 0xBC:
		return "Create primitive array"
	case 0xBD:
		return "Create reference array"
	case 0xBE:
		return "Get array length"

	default:
		return getOpcodeName(opcode)
	}
}

// getConstantPoolValue returns a readable representation of a constant pool entry
func getConstantPoolValue(frame *runtime.Frame, index uint16) string {
	cp := frame.Class.ConstantPool
	if int(index) >= len(cp) || cp[index] == nil {
		return "?"
	}

	entry := cp[index]
	switch c := entry.(type) {
	case *classfile.ConstantIntegerInfo:
		return fmt.Sprintf("int(%d)", c.Value)
	case *classfile.ConstantFloatInfo:
		return fmt.Sprintf("float(0x%X)", c.Value)
	case *classfile.ConstantLongInfo:
		return fmt.Sprintf("long(%d)", c.Value)
	case *classfile.ConstantDoubleInfo:
		return fmt.Sprintf("double(0x%X)", c.Value)
	case *classfile.ConstantStringInfo:
		str := cp.GetUtf8(c.StringIndex)
		if len(str) > 20 {
			return fmt.Sprintf("\"%s...\"", str[:17])
		}
		return fmt.Sprintf("\"%s\"", str)
	case *classfile.ConstantClassInfo:
		return fmt.Sprintf("class(%s)", cp.GetUtf8(c.NameIndex))
	case *classfile.ConstantMethodrefInfo:
		className := cp.GetClassName(c.ClassIndex)
		methodName, desc := cp.GetNameAndType(c.NameAndTypeIndex)
		return fmt.Sprintf("%s.%s%s", className, methodName, desc)
	case *classfile.ConstantFieldrefInfo:
		className := cp.GetClassName(c.ClassIndex)
		fieldName, _ := cp.GetNameAndType(c.NameAndTypeIndex)
		return fmt.Sprintf("%s.%s", className, fieldName)
	default:
		return fmt.Sprintf("cp[%d]", index)
	}
}

// PrintConstantPool prints the constant pool in a readable format
func PrintConstantPool(cp classfile.ConstantPool) {
	fmt.Println("┌─ Constant Pool ─────────────────────────────────────────────")
	for i := 1; i < len(cp); i++ {
		entry := cp[i]
		if entry == nil {
			continue // Skip empty slots (after Long/Double)
		}

		var desc string
		switch c := entry.(type) {
		case *classfile.ConstantUtf8Info:
			val := c.Value
			if len(val) > 40 {
				val = val[:37] + "..."
			}
			desc = fmt.Sprintf("Utf8          \"%s\"", val)
		case *classfile.ConstantIntegerInfo:
			desc = fmt.Sprintf("Integer       %d", c.Value)
		case *classfile.ConstantFloatInfo:
			desc = fmt.Sprintf("Float         0x%X", c.Value)
		case *classfile.ConstantLongInfo:
			desc = fmt.Sprintf("Long          %d", c.Value)
		case *classfile.ConstantDoubleInfo:
			desc = fmt.Sprintf("Double        0x%X", c.Value)
		case *classfile.ConstantClassInfo:
			name := cp.GetUtf8(c.NameIndex)
			desc = fmt.Sprintf("Class         #%d → %s", c.NameIndex, name)
		case *classfile.ConstantStringInfo:
			str := cp.GetUtf8(c.StringIndex)
			if len(str) > 30 {
				str = str[:27] + "..."
			}
			desc = fmt.Sprintf("String        #%d → \"%s\"", c.StringIndex, str)
		case *classfile.ConstantFieldrefInfo:
			className := cp.GetClassName(c.ClassIndex)
			fieldName, fieldDesc := cp.GetNameAndType(c.NameAndTypeIndex)
			desc = fmt.Sprintf("Fieldref      %s.%s:%s", className, fieldName, fieldDesc)
		case *classfile.ConstantMethodrefInfo:
			className := cp.GetClassName(c.ClassIndex)
			methodName, methodDesc := cp.GetNameAndType(c.NameAndTypeIndex)
			desc = fmt.Sprintf("Methodref     %s.%s%s", className, methodName, methodDesc)
		case *classfile.ConstantInterfaceMethodrefInfo:
			className := cp.GetClassName(c.ClassIndex)
			methodName, methodDesc := cp.GetNameAndType(c.NameAndTypeIndex)
			desc = fmt.Sprintf("InterfaceRef  %s.%s%s", className, methodName, methodDesc)
		case *classfile.ConstantNameAndTypeInfo:
			name := cp.GetUtf8(c.NameIndex)
			typeDesc := cp.GetUtf8(c.DescriptorIndex)
			desc = fmt.Sprintf("NameAndType   %s:%s", name, typeDesc)
		default:
			desc = fmt.Sprintf("%T", entry)
		}

		fmt.Printf("│ #%-3d  %s\n", i, desc)
	}
	fmt.Println("└─────────────────────────────────────────────────────────────")
	fmt.Println()
}

// formatRef formats a reference value for display
func formatRef(ref any) string {
	if ref == nil {
		return "null"
	}
	switch v := ref.(type) {
	case *runtime.Object:
		return fmt.Sprintf("Object<%s>@%p", v.ClassName(), v)
	case *runtime.Array:
		return fmt.Sprintf("Array[%d]@%p", v.Length, v)
	case string:
		if len(v) > 20 {
			return fmt.Sprintf("\"%s...\"", v[:20])
		}
		return fmt.Sprintf("\"%s\"", v)
	default:
		return fmt.Sprintf("%T", v)
	}
}

// getOpcodeName returns a human-readable name for an opcode
func getOpcodeName(opcode uint8) string {
	names := map[uint8]string{
		0x00: "nop", 0x01: "aconst_null",
		0x02: "iconst_m1", 0x03: "iconst_0", 0x04: "iconst_1", 0x05: "iconst_2",
		0x06: "iconst_3", 0x07: "iconst_4", 0x08: "iconst_5",
		0x09: "lconst_0", 0x0A: "lconst_1",
		0x10: "bipush", 0x11: "sipush",
		0x12: "ldc", 0x13: "ldc_w", 0x14: "ldc2_w",
		0x15: "iload", 0x16: "lload", 0x19: "aload",
		0x1A: "iload_0", 0x1B: "iload_1", 0x1C: "iload_2", 0x1D: "iload_3",
		0x1E: "lload_0", 0x1F: "lload_1", 0x20: "lload_2", 0x21: "lload_3",
		0x2A: "aload_0", 0x2B: "aload_1", 0x2C: "aload_2", 0x2D: "aload_3",
		0x2E: "iaload", 0x2F: "laload", 0x30: "faload", 0x31: "daload",
		0x32: "aaload", 0x33: "baload", 0x34: "caload", 0x35: "saload",
		0x36: "istore", 0x37: "lstore", 0x3A: "astore",
		0x3B: "istore_0", 0x3C: "istore_1", 0x3D: "istore_2", 0x3E: "istore_3",
		0x3F: "lstore_0", 0x40: "lstore_1", 0x41: "lstore_2", 0x42: "lstore_3",
		0x4B: "astore_0", 0x4C: "astore_1", 0x4D: "astore_2", 0x4E: "astore_3",
		0x4F: "iastore", 0x50: "lastore", 0x51: "fastore", 0x52: "dastore",
		0x53: "aastore", 0x54: "bastore", 0x55: "castore", 0x56: "sastore",
		0x57: "pop", 0x58: "pop2", 0x59: "dup", 0x5A: "dup_x1",
		0x5B: "dup_x2", 0x5C: "dup2", 0x5F: "swap",
		0x60: "iadd", 0x61: "ladd", 0x64: "isub", 0x65: "lsub",
		0x68: "imul", 0x69: "lmul", 0x6C: "idiv", 0x6D: "ldiv",
		0x70: "irem", 0x71: "lrem", 0x74: "ineg", 0x75: "lneg",
		0x78: "ishl", 0x79: "lshl", 0x7A: "ishr", 0x7B: "lshr",
		0x7C: "iushr", 0x7D: "lushr",
		0x7E: "iand", 0x7F: "land", 0x80: "ior", 0x81: "lor",
		0x82: "ixor", 0x83: "lxor", 0x84: "iinc",
		0x85: "i2l", 0x86: "i2f", 0x87: "i2d", 0x88: "l2i",
		0x94: "lcmp",
		0x99: "ifeq", 0x9A: "ifne", 0x9B: "iflt", 0x9C: "ifge",
		0x9D: "ifgt", 0x9E: "ifle",
		0x9F: "if_icmpeq", 0xA0: "if_icmpne", 0xA1: "if_icmplt",
		0xA2: "if_icmpge", 0xA3: "if_icmpgt", 0xA4: "if_icmple",
		0xA5: "if_acmpeq", 0xA6: "if_acmpne",
		0xA7: "goto", 0xA8: "jsr", 0xA9: "ret",
		0xAA: "tableswitch", 0xAB: "lookupswitch",
		0xAC: "ireturn", 0xAD: "lreturn", 0xAE: "freturn",
		0xAF: "dreturn", 0xB0: "areturn", 0xB1: "return",
		0xB2: "getstatic", 0xB3: "putstatic",
		0xB4: "getfield", 0xB5: "putfield",
		0xB6: "invokevirtual", 0xB7: "invokespecial",
		0xB8: "invokestatic", 0xB9: "invokeinterface", 0xBA: "invokedynamic",
		0xBB: "new", 0xBC: "newarray", 0xBD: "anewarray",
		0xBE: "arraylength", 0xBF: "athrow",
		0xC0: "checkcast", 0xC1: "instanceof",
		0xC2: "monitorenter", 0xC3: "monitorexit",
		0xC6: "ifnull", 0xC7: "ifnonnull", 0xC8: "goto_w",
	}
	if name, ok := names[opcode]; ok {
		return name
	}
	return "unknown"
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

// trackAlloc registers an object with the heap for GC tracking
func (i *Interpreter) trackAlloc(obj any) {
	if jvm := i.thread.JVM(); jvm != nil {
		jvm.GetHeap().Alloc(obj)
	}
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
