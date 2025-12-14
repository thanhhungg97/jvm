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

## Phase 1: Arrays & Strings

```go
// Add to interpreter
case NEWARRAY:     // Create primitive array
case ANEWARRAY:    // Create reference array
case ARRAYLENGTH:  // Get array length
case IALOAD:       // Load int from array
case IASTORE:      // Store int to array
case AALOAD:       // Load reference from array
case AASTORE:      // Store reference to array
case CALOAD:       // Load char from array (for strings)
case CASTORE:      // Store char to array
```

Benefits: Process command line args, use String methods

## Phase 2: Object Support

```go
// Add to runtime
type Object struct {
    Class  *classfile.ClassFile
    Fields map[string]interface{}
}

// Add to interpreter
case NEW:           // Create new object (already partial)
case INVOKESPECIAL: // Call constructor <init>
case INVOKEVIRTUAL: // Call instance method
case GETFIELD:      // Get object field
case PUTFIELD:      // Set object field
```

Benefits: Create objects, call instance methods

## Phase 3: Exception Handling

```go
// Add to frame
type Frame struct {
    // ... existing fields
    ExceptionHandlers []*ExceptionHandler
}

// Add to interpreter
case ATHROW:        // Throw exception
// Handle exception table lookups
```

Benefits: try/catch/finally, proper error handling

## Phase 4: Inheritance & Interfaces

```go
// Add to runtime
func (t *Thread) ResolveMethod(class, name, descriptor string) *MethodInfo {
    // Walk up class hierarchy
    // Check interfaces
}

// Add to interpreter
case CHECKCAST:     // Type checking
case INSTANCEOF:    // Type testing
```

Benefits: Polymorphism, interface implementations

## Phase 5: Native Methods

```go
// Create native method registry
var nativeMethods = map[string]func(frame *Frame){
    "java/lang/System.currentTimeMillis": nativeCurrentTimeMillis,
    "java/io/FileInputStream.read0":      nativeFileRead,
    "java/net/Socket.connect0":           nativeSocketConnect,
}
```

Benefits: File I/O, networking, system calls

## Phase 6: Threading

```go
// Add to runtime
type JVM struct {
    Threads []*Thread
    Heap    *Heap
    mutex   sync.Mutex
}

// Add to interpreter
case MONITORENTER:  // synchronized block enter
case MONITOREXIT:   // synchronized block exit
```

Benefits: Concurrent code, parallel execution

## Phase 7: Garbage Collection

```go
// Add to runtime
type Heap struct {
    objects   []*Object
    allocator *Allocator
}

func (h *Heap) GC() {
    // Mark reachable objects
    // Sweep unreachable objects
}
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

