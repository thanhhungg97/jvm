package runtime

import (
	"testing"
)

// ==================== Array Tests ====================

func TestNewPrimitiveArray(t *testing.T) {
	tests := []struct {
		name   string
		atype  ArrayType
		length int32
	}{
		{"int array", ArrayTypeInt, 10},
		{"long array", ArrayTypeLong, 5},
		{"byte array", ArrayTypeByte, 100},
		{"char array", ArrayTypeChar, 50},
		{"boolean array", ArrayTypeBoolean, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := NewPrimitiveArray(tt.atype, tt.length)
			if arr == nil {
				t.Fatal("NewPrimitiveArray returned nil")
				return // unreachable but satisfies staticcheck
			}
			if arr.Length != tt.length {
				t.Errorf("Length = %d, want %d", arr.Length, tt.length)
			}
			if arr.Type != tt.atype {
				t.Errorf("Type = %d, want %d", arr.Type, tt.atype)
			}
			if arr.IsRefArray() {
				t.Error("Primitive array should not be a reference array")
			}
		})
	}
}

func TestArrayIntOperations(t *testing.T) {
	arr := NewPrimitiveArray(ArrayTypeInt, 5)

	// Test set and get
	arr.SetInt(0, 100)
	arr.SetInt(2, 200)
	arr.SetInt(4, 300)

	if arr.GetInt(0) != 100 {
		t.Errorf("GetInt(0) = %d, want 100", arr.GetInt(0))
	}
	if arr.GetInt(2) != 200 {
		t.Errorf("GetInt(2) = %d, want 200", arr.GetInt(2))
	}
	if arr.GetInt(4) != 300 {
		t.Errorf("GetInt(4) = %d, want 300", arr.GetInt(4))
	}
	// Unset index should be zero
	if arr.GetInt(1) != 0 {
		t.Errorf("GetInt(1) = %d, want 0", arr.GetInt(1))
	}
}

func TestReferenceArray(t *testing.T) {
	arr := NewReferenceArray("java/lang/String", 3)

	if !arr.IsRefArray() {
		t.Error("Reference array should be a reference array")
	}
	if arr.Length != 3 {
		t.Errorf("Length = %d, want 3", arr.Length)
	}
	if arr.ClassName != "java/lang/String" {
		t.Errorf("ClassName = %s, want java/lang/String", arr.ClassName)
	}

	// Test set and get
	arr.SetRef(0, "Hello")
	arr.SetRef(1, "World")

	if arr.GetRef(0) != "Hello" {
		t.Errorf("GetRef(0) = %v, want Hello", arr.GetRef(0))
	}
	if arr.GetRef(1) != "World" {
		t.Errorf("GetRef(1) = %v, want World", arr.GetRef(1))
	}
	if arr.GetRef(2) != nil {
		t.Errorf("GetRef(2) = %v, want nil", arr.GetRef(2))
	}
}

// ==================== OperandStack Tests ====================

func TestOperandStackInt(t *testing.T) {
	stack := NewOperandStack(10)

	stack.PushInt(42)
	stack.PushInt(-100)
	stack.PushInt(999)

	if stack.Size() != 3 {
		t.Errorf("Size() = %d, want 3", stack.Size())
	}

	if val := stack.PopInt(); val != 999 {
		t.Errorf("PopInt() = %d, want 999", val)
	}
	if val := stack.PopInt(); val != -100 {
		t.Errorf("PopInt() = %d, want -100", val)
	}
	if val := stack.PopInt(); val != 42 {
		t.Errorf("PopInt() = %d, want 42", val)
	}

	if !stack.IsEmpty() {
		t.Error("Stack should be empty")
	}
}

func TestOperandStackLong(t *testing.T) {
	stack := NewOperandStack(10)

	stack.PushLong(1234567890123)
	stack.PushLong(-9876543210)

	if val := stack.PopLong(); val != -9876543210 {
		t.Errorf("PopLong() = %d, want -9876543210", val)
	}
	if val := stack.PopLong(); val != 1234567890123 {
		t.Errorf("PopLong() = %d, want 1234567890123", val)
	}
}

func TestOperandStackRef(t *testing.T) {
	stack := NewOperandStack(10)

	stack.PushRef("Hello")
	stack.PushRef(nil)
	arr := NewPrimitiveArray(ArrayTypeInt, 5)
	stack.PushRef(arr)

	if ref := stack.PopRef(); ref != arr {
		t.Errorf("PopRef() = %v, want array", ref)
	}
	if ref := stack.PopRef(); ref != nil {
		t.Errorf("PopRef() = %v, want nil", ref)
	}
	if ref := stack.PopRef(); ref != "Hello" {
		t.Errorf("PopRef() = %v, want Hello", ref)
	}
}

func TestOperandStackDup(t *testing.T) {
	stack := NewOperandStack(10)

	stack.PushInt(42)
	stack.Dup()

	if stack.Size() != 2 {
		t.Errorf("Size() = %d, want 2", stack.Size())
	}
	if val := stack.PopInt(); val != 42 {
		t.Errorf("PopInt() = %d, want 42", val)
	}
	if val := stack.PopInt(); val != 42 {
		t.Errorf("PopInt() = %d, want 42", val)
	}
}

func TestOperandStackSwap(t *testing.T) {
	stack := NewOperandStack(10)

	stack.PushInt(1)
	stack.PushInt(2)
	stack.Swap()

	if val := stack.PopInt(); val != 1 {
		t.Errorf("PopInt() = %d, want 1", val)
	}
	if val := stack.PopInt(); val != 2 {
		t.Errorf("PopInt() = %d, want 2", val)
	}
}

// ==================== LocalVars Tests ====================

func TestLocalVars(t *testing.T) {
	locals := NewLocalVars(10)

	// Test int
	locals.SetInt(0, 100)
	if val := locals.GetInt(0); val != 100 {
		t.Errorf("GetInt(0) = %d, want 100", val)
	}

	// Test long
	locals.SetLong(1, 9876543210)
	if val := locals.GetLong(1); val != 9876543210 {
		t.Errorf("GetLong(1) = %d, want 9876543210", val)
	}

	// Test ref
	locals.SetRef(3, "test")
	if val := locals.GetRef(3); val != "test" {
		t.Errorf("GetRef(3) = %v, want test", val)
	}
}

// ==================== Heap Tests ====================

func TestHeapAlloc(t *testing.T) {
	heap := NewHeap()
	heap.SetGCEnabled(false) // Disable GC for testing

	id1 := heap.Alloc("object1")
	id2 := heap.Alloc("object2")
	id3 := heap.Alloc("object3")

	if id1 == id2 || id2 == id3 || id1 == id3 {
		t.Error("Alloc should return unique IDs")
	}

	if obj := heap.Get(id1); obj != "object1" {
		t.Errorf("Get(id1) = %v, want object1", obj)
	}
	if obj := heap.Get(id2); obj != "object2" {
		t.Errorf("Get(id2) = %v, want object2", obj)
	}

	stats := heap.Stats()
	if stats.AllocCount != 3 {
		t.Errorf("AllocCount = %d, want 3", stats.AllocCount)
	}
}

func TestHeapFree(t *testing.T) {
	heap := NewHeap()
	heap.SetGCEnabled(false)

	id := heap.Alloc("test")
	heap.Free(id)

	if obj := heap.Get(id); obj != nil {
		t.Errorf("Get after Free = %v, want nil", obj)
	}

	stats := heap.Stats()
	if stats.FreeCount != 1 {
		t.Errorf("FreeCount = %d, want 1", stats.FreeCount)
	}
}

// ==================== Monitor Tests ====================

func TestMonitorEnterExit(t *testing.T) {
	jvm := NewJVM()
	thread := jvm.CreateThread()

	obj := "test-object"
	monitor := jvm.GetOrCreateMonitor(obj)

	// Enter should succeed
	monitor.Enter(thread)

	// Exit should succeed
	err := monitor.Exit(thread)
	if err != nil {
		t.Errorf("Exit failed: %v", err)
	}
}

func TestMonitorReentrant(t *testing.T) {
	jvm := NewJVM()
	thread := jvm.CreateThread()

	obj := "test-object"
	monitor := jvm.GetOrCreateMonitor(obj)

	// Enter multiple times (reentrant)
	monitor.Enter(thread)
	monitor.Enter(thread)
	monitor.Enter(thread)

	// Should need 3 exits
	if err := monitor.Exit(thread); err != nil {
		t.Errorf("Exit 1 failed: %v", err)
	}
	if err := monitor.Exit(thread); err != nil {
		t.Errorf("Exit 2 failed: %v", err)
	}
	if err := monitor.Exit(thread); err != nil {
		t.Errorf("Exit 3 failed: %v", err)
	}
}

// ==================== Native Registry Tests ====================

func TestNativeRegistry(t *testing.T) {
	// Test that built-in natives are registered
	natives := []struct {
		class      string
		method     string
		descriptor string
	}{
		{"java/lang/System", "currentTimeMillis", "()J"},
		{"java/lang/System", "nanoTime", "()J"},
		{"java/lang/Math", "abs", "(I)I"},
		{"java/lang/Thread", "sleep", "(J)V"},
	}

	for _, n := range natives {
		t.Run(n.class+"."+n.method, func(t *testing.T) {
			method := Natives.Lookup(n.class, n.method, n.descriptor)
			if method == nil {
				t.Errorf("Native method %s.%s%s not found", n.class, n.method, n.descriptor)
			}
		})
	}
}

// ==================== JVM Tests ====================

func TestJVMCreateThread(t *testing.T) {
	jvm := NewJVM()

	t1 := jvm.CreateThread()
	t2 := jvm.CreateThread()

	if t1 == nil || t2 == nil {
		t.Fatal("CreateThread returned nil")
	}

	if t1.ID() == t2.ID() {
		t.Error("Threads should have unique IDs")
	}

	if jvm.GetMainThread() != t1 {
		t.Error("First thread should be main thread")
	}
}

func TestJVMIsRunning(t *testing.T) {
	jvm := NewJVM()

	if !jvm.IsRunning() {
		t.Error("New JVM should be running")
	}

	jvm.Shutdown()

	if jvm.IsRunning() {
		t.Error("JVM should not be running after shutdown")
	}
}
