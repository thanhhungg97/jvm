package runtime

import (
	"fmt"
	"time"
	"unsafe"
)

// NativeMethod represents a native method implementation
type NativeMethod func(frame *Frame) error

// NativeRegistry holds all registered native methods
type NativeRegistry struct {
	methods map[string]NativeMethod
}

// Global native registry
var Natives = NewNativeRegistry()

// NewNativeRegistry creates a new native method registry
func NewNativeRegistry() *NativeRegistry {
	nr := &NativeRegistry{
		methods: make(map[string]NativeMethod),
	}
	nr.registerBuiltins()
	return nr
}

// Register registers a native method
func (nr *NativeRegistry) Register(className, methodName, descriptor string, method NativeMethod) {
	key := className + "." + methodName + descriptor
	nr.methods[key] = method
}

// Lookup finds a native method
func (nr *NativeRegistry) Lookup(className, methodName, descriptor string) NativeMethod {
	key := className + "." + methodName + descriptor
	method := nr.methods[key]
	return method
}

// ListAll returns all registered native method keys (for debugging)
func (nr *NativeRegistry) ListAll() []string {
	keys := make([]string, 0, len(nr.methods))
	for k := range nr.methods {
		keys = append(keys, k)
	}
	return keys
}

// Count returns the number of registered native methods
func (nr *NativeRegistry) Count() int {
	return len(nr.methods)
}

// registerBuiltins registers all built-in native methods
func (nr *NativeRegistry) registerBuiltins() {
	// System class natives
	nr.Register("java/lang/System", "currentTimeMillis", "()J", nativeCurrentTimeMillis)
	nr.Register("java/lang/System", "nanoTime", "()J", nativeNanoTime)
	nr.Register("java/lang/System", "arraycopy", "(Ljava/lang/Object;ILjava/lang/Object;II)V", nativeArraycopy)
	nr.Register("java/lang/System", "identityHashCode", "(Ljava/lang/Object;)I", nativeIdentityHashCode)

	// Object class natives
	nr.Register("java/lang/Object", "hashCode", "()I", nativeObjectHashCode)
	nr.Register("java/lang/Object", "getClass", "()Ljava/lang/Class;", nativeObjectGetClass)

	// Class class natives
	nr.Register("java/lang/Class", "getName", "()Ljava/lang/String;", nativeClassGetName)
	nr.Register("java/lang/Class", "isPrimitive", "()Z", nativeClassIsPrimitive)

	// Thread class natives
	nr.Register("java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", nativeThreadCurrentThread)
	nr.Register("java/lang/Thread", "sleep", "(J)V", nativeThreadSleep)

	// Math class natives
	nr.Register("java/lang/Math", "sqrt", "(D)D", nativeMathSqrt)
	nr.Register("java/lang/Math", "abs", "(I)I", nativeMathAbsInt)
	nr.Register("java/lang/Math", "abs", "(J)J", nativeMathAbsLong)
	nr.Register("java/lang/Math", "max", "(II)I", nativeMathMaxInt)
	nr.Register("java/lang/Math", "min", "(II)I", nativeMathMinInt)

	// String class natives
	nr.Register("java/lang/String", "intern", "()Ljava/lang/String;", nativeStringIntern)

	// Float/Double natives
	nr.Register("java/lang/Float", "floatToRawIntBits", "(F)I", nativeFloatToRawIntBits)
	nr.Register("java/lang/Double", "doubleToRawLongBits", "(D)J", nativeDoubleToRawLongBits)

	// Runtime natives
	nr.Register("java/lang/Runtime", "availableProcessors", "()I", nativeAvailableProcessors)
	nr.Register("java/lang/Runtime", "freeMemory", "()J", nativeFreeMemory)
	nr.Register("java/lang/Runtime", "totalMemory", "()J", nativeTotalMemory)
	nr.Register("java/lang/Runtime", "maxMemory", "()J", nativeMaxMemory)
	nr.Register("java/lang/Runtime", "gc", "()V", nativeGC)
}

// =============== System natives ===============

func nativeCurrentTimeMillis(frame *Frame) error {
	millis := time.Now().UnixMilli()
	frame.OperandStack.PushLong(millis)
	return nil
}

func nativeNanoTime(frame *Frame) error {
	nanos := time.Now().UnixNano()
	frame.OperandStack.PushLong(nanos)
	return nil
}

func nativeArraycopy(frame *Frame) error {
	stack := frame.OperandStack
	length := stack.PopInt()
	destPos := stack.PopInt()
	destRef := stack.PopRef()
	srcPos := stack.PopInt()
	srcRef := stack.PopRef()

	if srcRef == nil || destRef == nil {
		return fmt.Errorf("NullPointerException: arraycopy with null array")
	}

	srcArr, srcOk := srcRef.(*Array)
	destArr, destOk := destRef.(*Array)
	if !srcOk || !destOk {
		return fmt.Errorf("ArrayStoreException: not arrays")
	}

	// Copy elements
	for i := int32(0); i < length; i++ {
		if srcArr.IsRefArray() {
			destArr.SetRef(destPos+i, srcArr.GetRef(srcPos+i))
		} else {
			destArr.SetInt(destPos+i, srcArr.GetInt(srcPos+i))
		}
	}

	return nil
}

func nativeIdentityHashCode(frame *Frame) error {
	obj := frame.OperandStack.PopRef()
	if obj == nil {
		frame.OperandStack.PushInt(0)
	} else {
		// Use pointer address as hash code
		frame.OperandStack.PushInt(int32(uintptr(fmt.Sprintf("%p", obj)[2:][0])))
	}
	return nil
}

// =============== Object natives ===============

func nativeObjectHashCode(frame *Frame) error {
	// 'this' is in local var 0
	obj := frame.LocalVars.GetRef(0)
	if obj == nil {
		frame.OperandStack.PushInt(0)
	} else {
		// Simple hash based on object address
		frame.OperandStack.PushInt(int32(time.Now().UnixNano() & 0x7FFFFFFF))
	}
	return nil
}

func nativeObjectGetClass(frame *Frame) error {
	obj := frame.LocalVars.GetRef(0)
	if o, ok := obj.(*Object); ok && o.Class != nil {
		// Return class name as a placeholder for Class object
		frame.OperandStack.PushRef("Class<" + o.Class.ClassName() + ">")
	} else {
		frame.OperandStack.PushRef(nil)
	}
	return nil
}

// =============== Class natives ===============

func nativeClassGetName(frame *Frame) error {
	classRef := frame.LocalVars.GetRef(0)
	if str, ok := classRef.(string); ok {
		// Extract name from "Class<name>"
		if len(str) > 6 && str[:6] == "Class<" && str[len(str)-1] == '>' {
			name := str[6 : len(str)-1]
			// Convert java/lang/Object to java.lang.Object
			result := ""
			for _, c := range name {
				if c == '/' {
					result += "."
				} else {
					result += string(c)
				}
			}
			frame.OperandStack.PushRef(result)
			return nil
		}
	}
	frame.OperandStack.PushRef("Unknown")
	return nil
}

func nativeClassIsPrimitive(frame *Frame) error {
	classRef := frame.LocalVars.GetRef(0)
	if str, ok := classRef.(string); ok {
		// Primitive class names
		primitives := map[string]bool{
			"Class<int>": true, "Class<long>": true, "Class<float>": true,
			"Class<double>": true, "Class<boolean>": true, "Class<byte>": true,
			"Class<char>": true, "Class<short>": true, "Class<void>": true,
		}
		if primitives[str] {
			frame.OperandStack.PushInt(1)
			return nil
		}
	}
	frame.OperandStack.PushInt(0)
	return nil
}

// =============== Thread natives ===============

func nativeThreadCurrentThread(frame *Frame) error {
	// Return a placeholder thread object
	frame.OperandStack.PushRef("Thread<main>")
	return nil
}

func nativeThreadSleep(frame *Frame) error {
	millis := frame.OperandStack.PopLong()
	time.Sleep(time.Duration(millis) * time.Millisecond)
	return nil
}

// =============== Math natives ===============

func nativeMathSqrt(frame *Frame) error {
	// Pop double, compute sqrt, push result
	val := frame.OperandStack.PopDouble()
	// Simple Newton's method for sqrt
	if val < 0 {
		frame.OperandStack.PushDouble(0) // NaN not supported
	} else if val == 0 {
		frame.OperandStack.PushDouble(0)
	} else {
		x := val
		for i := 0; i < 20; i++ {
			x = (x + val/x) / 2
		}
		frame.OperandStack.PushDouble(x)
	}
	return nil
}

func nativeMathAbsInt(frame *Frame) error {
	val := frame.OperandStack.PopInt()
	if val < 0 {
		val = -val
	}
	frame.OperandStack.PushInt(val)
	return nil
}

func nativeMathAbsLong(frame *Frame) error {
	val := frame.OperandStack.PopLong()
	if val < 0 {
		val = -val
	}
	frame.OperandStack.PushLong(val)
	return nil
}

func nativeMathMaxInt(frame *Frame) error {
	b := frame.OperandStack.PopInt()
	a := frame.OperandStack.PopInt()
	if a > b {
		frame.OperandStack.PushInt(a)
	} else {
		frame.OperandStack.PushInt(b)
	}
	return nil
}

func nativeMathMinInt(frame *Frame) error {
	b := frame.OperandStack.PopInt()
	a := frame.OperandStack.PopInt()
	if a < b {
		frame.OperandStack.PushInt(a)
	} else {
		frame.OperandStack.PushInt(b)
	}
	return nil
}

// =============== String natives ===============

// String pool for interning
var stringPool = make(map[string]string)

func nativeStringIntern(frame *Frame) error {
	str := frame.LocalVars.GetRef(0)
	if s, ok := str.(string); ok {
		if interned, exists := stringPool[s]; exists {
			frame.OperandStack.PushRef(interned)
		} else {
			stringPool[s] = s
			frame.OperandStack.PushRef(s)
		}
	} else {
		frame.OperandStack.PushRef(str)
	}
	return nil
}

// =============== Float/Double natives ===============

func nativeFloatToRawIntBits(frame *Frame) error {
	f := frame.OperandStack.PopFloat()
	bits := *(*int32)(unsafe.Pointer(&f))
	frame.OperandStack.PushInt(bits)
	return nil
}

func nativeDoubleToRawLongBits(frame *Frame) error {
	d := frame.OperandStack.PopDouble()
	bits := *(*int64)(unsafe.Pointer(&d))
	frame.OperandStack.PushLong(bits)
	return nil
}

// =============== Runtime natives ===============

func nativeAvailableProcessors(frame *Frame) error {
	frame.OperandStack.PushInt(1) // Single-threaded JVM
	return nil
}

func nativeFreeMemory(frame *Frame) error {
	frame.OperandStack.PushLong(100 * 1024 * 1024) // 100MB placeholder
	return nil
}

func nativeTotalMemory(frame *Frame) error {
	frame.OperandStack.PushLong(256 * 1024 * 1024) // 256MB placeholder
	return nil
}

func nativeMaxMemory(frame *Frame) error {
	frame.OperandStack.PushLong(512 * 1024 * 1024) // 512MB placeholder
	return nil
}

func nativeGC(frame *Frame) error {
	// Placeholder - actual GC will be implemented in Phase 7
	return nil
}

// =============== File I/O natives ===============

// RegisterFileNatives adds file I/O native methods
func init() {
	// These would be registered if we had a proper FileInputStream/FileOutputStream
	// For now, we provide simple console I/O through System.out which is already handled
}
