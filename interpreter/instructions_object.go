package interpreter

import (
	"fmt"
	"simplejvm/classfile"
	"simplejvm/runtime"
)

// executeObjectInstruction handles object-related instructions
func (i *Interpreter) executeObjectInstruction(frame *runtime.Frame, opcode uint8) (bool, error) {
	stack := frame.OperandStack
	cp := frame.Class.ConstantPool

	switch opcode {
	// Get static field
	case GETSTATIC:
		index := frame.ReadU2()
		fieldRef := cp[index].(*classfile.ConstantFieldrefInfo)
		className := cp.GetClassName(fieldRef.ClassIndex)
		fieldName, descriptor := cp.GetNameAndType(fieldRef.NameAndTypeIndex)

		if className == "java/lang/System" && fieldName == "out" {
			stack.PushRef("System.out")
		} else {
			key := className + "." + fieldName
			if val, ok := i.staticFields[key]; ok {
				switch descriptor[0] {
				case 'B', 'C', 'I', 'S', 'Z':
					stack.PushInt(val.(int32))
				case 'J':
					stack.PushLong(val.(int64))
				case 'L', '[':
					stack.PushRef(val)
				default:
					stack.PushRef(val)
				}
			} else {
				switch descriptor[0] {
				case 'B', 'C', 'I', 'S', 'Z':
					stack.PushInt(0)
				case 'J':
					stack.PushLong(0)
				default:
					stack.PushRef(nil)
				}
			}
		}

	// Put static field
	case PUTSTATIC:
		index := frame.ReadU2()
		fieldRef := cp[index].(*classfile.ConstantFieldrefInfo)
		className := cp.GetClassName(fieldRef.ClassIndex)
		fieldName, descriptor := cp.GetNameAndType(fieldRef.NameAndTypeIndex)
		key := className + "." + fieldName

		switch descriptor[0] {
		case 'B', 'C', 'I', 'S', 'Z':
			i.staticFields[key] = stack.PopInt()
		case 'J':
			i.staticFields[key] = stack.PopLong()
		case 'L', '[':
			i.staticFields[key] = stack.PopRef()
		default:
			i.staticFields[key] = stack.PopRef()
		}

	// Get instance field
	case GETFIELD:
		index := frame.ReadU2()
		fieldRef := cp[index].(*classfile.ConstantFieldrefInfo)
		_, descriptor := cp.GetNameAndType(fieldRef.NameAndTypeIndex)
		fieldName, _ := cp.GetNameAndType(fieldRef.NameAndTypeIndex)

		objRef := stack.PopRef()
		if objRef == nil {
			return true, fmt.Errorf("NullPointerException: getfield on null object")
		}

		obj, ok := objRef.(*runtime.Object)
		if !ok {
			return true, fmt.Errorf("getfield: not an Object: %T", objRef)
		}

		switch descriptor[0] {
		case 'B', 'C', 'I', 'S', 'Z':
			stack.PushInt(obj.GetFieldInt(fieldName))
		case 'J':
			stack.PushLong(obj.GetFieldLong(fieldName))
		case 'L', '[':
			stack.PushRef(obj.GetFieldRef(fieldName))
		default:
			return true, fmt.Errorf("getfield: unsupported field type: %s", descriptor)
		}

	// Put instance field
	case PUTFIELD:
		index := frame.ReadU2()
		fieldRef := cp[index].(*classfile.ConstantFieldrefInfo)
		_, descriptor := cp.GetNameAndType(fieldRef.NameAndTypeIndex)
		fieldName, _ := cp.GetNameAndType(fieldRef.NameAndTypeIndex)

		switch descriptor[0] {
		case 'B', 'C', 'I', 'S', 'Z':
			val := stack.PopInt()
			objRef := stack.PopRef()
			if objRef == nil {
				return true, fmt.Errorf("NullPointerException: putfield on null object")
			}
			obj := objRef.(*runtime.Object)
			obj.SetFieldInt(fieldName, val)
		case 'J':
			val := stack.PopLong()
			objRef := stack.PopRef()
			if objRef == nil {
				return true, fmt.Errorf("NullPointerException: putfield on null object")
			}
			obj := objRef.(*runtime.Object)
			obj.SetFieldLong(fieldName, val)
		case 'L', '[':
			val := stack.PopRef()
			objRef := stack.PopRef()
			if objRef == nil {
				return true, fmt.Errorf("NullPointerException: putfield on null object")
			}
			obj := objRef.(*runtime.Object)
			obj.SetFieldRef(fieldName, val)
		default:
			return true, fmt.Errorf("putfield: unsupported field type: %s", descriptor)
		}

	// New object
	case NEW:
		index := frame.ReadU2()
		className := cp.GetClassName(index)
		cf := i.thread.GetClass(className)
		if cf == nil && className == frame.Class.ClassName() {
			cf = frame.Class
		}
		// Try to lazy-load the class if not found
		if cf == nil {
			cf = i.tryLoadClass(className, frame.Class)
		}
		var obj interface{}
		if cf != nil {
			obj = runtime.NewObject(cf)
			i.trackAlloc(obj) // Track on heap for GC
		} else {
			// Fallback to placeholder for system classes we don't support
			obj = fmt.Sprintf("Object<%s>", className)
		}
		stack.PushRef(obj)

	// Type checking
	case CHECKCAST:
		index := frame.ReadU2()
		targetClassName := cp.GetClassName(index)
		objRef := stack.PopRef()

		if objRef == nil {
			stack.PushRef(nil)
		} else {
			objClassName := getObjectClassName(objRef)
			if objClassName == targetClassName || targetClassName == "java/lang/Object" {
				stack.PushRef(objRef)
			} else {
				stack.PushRef(objRef) // Lenient for now
			}
		}

	// Instance of check
	case INSTANCEOF:
		index := frame.ReadU2()
		targetClassName := cp.GetClassName(index)
		objRef := stack.PopRef()

		if objRef == nil {
			stack.PushInt(0)
		} else {
			objClassName := getObjectClassName(objRef)
			if objClassName == targetClassName || targetClassName == "java/lang/Object" {
				stack.PushInt(1)
			} else {
				stack.PushInt(0)
			}
		}

	// Synchronization - monitor enter
	case MONITORENTER:
		objRef := stack.PopRef()
		if objRef == nil {
			return true, fmt.Errorf("NullPointerException: monitorenter on null")
		}
		if jvm := i.thread.JVM(); jvm != nil {
			monitor := jvm.GetOrCreateMonitor(objRef)
			monitor.Enter(i.thread)
		}

	// Synchronization - monitor exit
	case MONITOREXIT:
		objRef := stack.PopRef()
		if objRef == nil {
			return true, fmt.Errorf("NullPointerException: monitorexit on null")
		}
		if jvm := i.thread.JVM(); jvm != nil {
			monitor := jvm.GetOrCreateMonitor(objRef)
			if err := monitor.Exit(i.thread); err != nil {
				return true, err
			}
		}

	// Throw exception
	case ATHROW:
		exRef := stack.PopRef()
		if exRef == nil {
			return true, fmt.Errorf("NullPointerException: cannot throw null")
		}
		exClassName := ""
		if obj, ok := exRef.(*runtime.Object); ok && obj.Class != nil {
			exClassName = obj.Class.ClassName()
		} else if str, ok := exRef.(string); ok {
			exClassName = extractClassName(str)
		}

		if err := i.handleException(exRef, exClassName); err != nil {
			return true, err
		}

	default:
		return false, nil
	}
	return true, nil
}

// getObjectClassName returns the class name of an object reference
func getObjectClassName(objRef any) string {
	if obj, ok := objRef.(*runtime.Object); ok && obj.Class != nil {
		return obj.Class.ClassName()
	} else if arr, ok := objRef.(*runtime.Array); ok {
		if arr.IsRefArray() {
			return "[L" + arr.ClassName + ";"
		}
		return "[" + string(arrayTypeChar(arr.Type))
	} else if str, ok := objRef.(string); ok {
		return extractClassName(str)
	}
	return ""
}
