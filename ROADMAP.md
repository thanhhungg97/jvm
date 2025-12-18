# SimpleJVM Roadmap

## Current Status ✓

- [x] Class file parsing (magic, constant pool, methods, code)
- [x] Bytecode interpreter loop
- [x] Operand stack and local variables
- [x] Integer/Long arithmetic
- [x] Conditional branches (if/else)
- [x] Loops (for/while via goto)
- [x] Static method calls
- [x] Method return values
- [x] System.out.println (int, long, boolean, String)
- [x] Bitwise operations
- [x] Recursion

## Phase 1: Arrays & Strings ✓

```go
// Implemented in interpreter
case NEWARRAY:     // Create primitive array ✓
case ANEWARRAY:    // Create reference array ✓
case ARRAYLENGTH:  // Get array length ✓
case IALOAD:       // Load int from array ✓
case IASTORE:      // Store int to array ✓
case AALOAD:       // Load reference from array ✓
case AASTORE:      // Store reference to array ✓
case CALOAD:       // Load char from array (for strings) ✓
case CASTORE:      // Store char to array ✓
case LALOAD:       // Load long from array ✓
case LASTORE:      // Store long to array ✓
case BALOAD:       // Load byte from array ✓
case BASTORE:      // Store byte to array ✓
case SALOAD:       // Load short from array ✓
case SASTORE:      // Store short to array ✓
case FALOAD:       // Load float from array ✓
case FASTORE:      // Store float to array ✓
case DALOAD:       // Load double from array ✓
case DASTORE:      // Store double to array ✓
```

Added: `runtime/array.go` with Array type supporting all primitive types and references.

Benefits: Process command line args, use String methods

## Phase 2: Object Support ✓

```go
// Added to runtime/object.go
type Object struct {
    Class      *classfile.ClassFile  ✓
    Fields     map[string]any        ✓
    FieldSlots map[string]int64      ✓
}

// Implemented in interpreter
case NEW:           // Create new object ✓
case INVOKESPECIAL: // Call constructor <init> ✓
case INVOKEVIRTUAL: // Call instance method ✓
case GETFIELD:      // Get object field ✓
case PUTFIELD:      // Set object field ✓
```

Benefits: Create objects, call instance methods

## Phase 3: Exception Handling ✓

```go
// Added to runtime/exception.go
type JavaException struct {
    Object    *Object  ✓
    ClassName string   ✓
    Message   string   ✓
}

func FindExceptionHandler(...) int ✓

// Implemented in interpreter
case ATHROW:        // Throw exception ✓
// Exception table lookups ✓
// Runtime exceptions (ArithmeticException, etc.) ✓
```

Benefits: try/catch/finally, proper error handling

## Phase 4: Inheritance & Interfaces ✓

```go
// Implemented in interpreter
case CHECKCAST:     // Type checking ✓
case INSTANCEOF:    // Type testing ✓
```

Note: Basic type checking implemented. Full inheritance hierarchy walking planned for future.

Benefits: Polymorphism, interface implementations

## Phase 5: Native Methods ✓

```go
// Added to runtime/native.go
type NativeRegistry struct { ... } ✓

// Implemented natives:
- System.currentTimeMillis ✓
- System.nanoTime ✓
- System.arraycopy ✓
- Math.abs, Math.max, Math.min ✓
- Thread.sleep ✓
- Object.hashCode ✓
- and more...
```

Benefits: File I/O, networking, system calls

## Phase 6: Threading ✓

```go
// Added to runtime/jvm.go
type JVM struct {
    mainThread  *Thread          ✓
    threads     []*Thread        ✓
    monitors    map[any]*Monitor ✓
    heap        *Heap            ✓
}

type Monitor struct { ... } ✓

// Implemented in interpreter
case MONITORENTER:  // synchronized block enter ✓
case MONITOREXIT:   // synchronized block exit ✓
```

Benefits: Concurrent code, parallel execution

## Phase 7: Garbage Collection ✓

```go
// Added to runtime/heap.go
type Heap struct {
    objects   map[uint64]any ✓
    gcEnabled bool           ✓
}

func (h *Heap) GC(roots []any) ✓  // Mark-sweep GC
func (h *Heap) Stats() HeapStats ✓
```

Benefits: Long-running applications, memory management

---

## HTTP Server Requirements

To run a simple HTTP server, you need at minimum:

| Phase | Feature | Why Needed |
|-------|---------|------------|
| 1 | Arrays | Byte buffers for I/O |
| 2 | Objects | Socket, InputStream objects |
| 3 | Exceptions | Handle IOException |
| 5 | Native Methods | Actual socket operations |

## Simpler Network Alternative

Instead of implementing full Java networking, you could:

1. **Create Go-based networking** exposed as native methods
2. **Use a simple protocol** (not full HTTP)
3. **Implement minimal socket operations**

```go
// Native networking bridge
var networkNatives = map[string]func(frame *Frame){
    "SimpleSocket.connect": func(f *Frame) {
        host := f.PopString()
        port := f.PopInt()
        conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
        f.PushRef(conn)
    },
    "SimpleSocket.read": func(f *Frame) {
        conn := f.PopRef().(net.Conn)
        buf := make([]byte, 1024)
        n, _ := conn.Read(buf)
        f.PushString(string(buf[:n]))
    },
}
```

This bypasses most JVM complexity while providing networking capability.



