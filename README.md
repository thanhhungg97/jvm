# SimpleJVM

A minimal JVM implementation in Go for educational purposes.

## Features

- Class file parsing (constant pool, methods, attributes)
- Bytecode interpreter with comprehensive instruction support
- Operand stack and local variables (primitives and references)
- Arrays (primitive and reference arrays)
- Object creation and instance method invocation
- Static and instance field access
- Constructor invocation
- Exception handling (try/catch/throw)
- Type checking (instanceof, checkcast)
- Native methods (System.currentTimeMillis, Math.*, Thread.sleep, etc.)
- Synchronized blocks (monitorenter/monitorexit)
- JVM structure with heap and monitor management
- Mark-sweep garbage collection (foundation)
- Console output via `System.out.println`

## Supported Instructions

- Constants: `iconst_*`, `lconst_*`, `aconst_null`, `bipush`, `sipush`, `ldc`, `ldc_w`, `ldc2_w`
- Loads: `iload`, `iload_*`, `lload`, `lload_*`, `aload`, `aload_*`
- Stores: `istore`, `istore_*`, `lstore`, `lstore_*`, `astore`, `astore_*`
- Arrays: `newarray`, `anewarray`, `arraylength`, `iaload`, `iastore`, `laload`, `lastore`, `aaload`, `aastore`, `baload`, `bastore`, `caload`, `castore`, `saload`, `sastore`, `faload`, `fastore`, `daload`, `dastore`
- Arithmetic: `iadd`, `isub`, `imul`, `idiv`, `irem`, `ineg`, `ladd`, `lsub`, `lmul`, `ldiv`, `lrem`, `lneg`, `iinc`
- Bitwise: `iand`, `ior`, `ixor`, `ishl`, `ishr`, `iushr`, `land`, `lor`, `lxor`
- Comparisons: `if_icmp*`, `ifeq`, `ifne`, `iflt`, `ifge`, `ifgt`, `ifle`, `lcmp`, `ifnull`, `ifnonnull`
- Conversions: `i2l`, `l2i`
- Stack: `pop`, `pop2`, `dup`, `swap`
- Control: `goto`, `goto_w`, `return`, `ireturn`, `lreturn`, `areturn`
- Objects: `new`, `getfield`, `putfield`, `getstatic`, `putstatic`
- Types: `checkcast`, `instanceof`
- Exceptions: `athrow`
- Synchronization: `monitorenter`, `monitorexit`
- Invocation: `invokestatic`, `invokevirtual`, `invokespecial`

## Usage

```bash
# Build
go build -o simplejvm

# Run a class file
./simplejvm HelloWorld.class

# Or compile and run a Java file
javac HelloWorld.java
./simplejvm HelloWorld.class
```

## Example

Create `HelloWorld.java`:
```java
public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, JVM!");
        System.out.println(add(10, 20));
    }
    
    public static int add(int a, int b) {
        return a + b;
    }
}
```

Compile and run:
```bash
javac HelloWorld.java
./simplejvm HelloWorld.class
```

## Architecture

```
simplejvm/
├── classfile/          # Class file parser
│   ├── classfile.go    # Class file structure
│   ├── constant_pool.go # Constant pool parsing
│   └── reader.go       # Binary reader
├── runtime/            # Runtime structures
│   ├── array.go        # Array types
│   ├── exception.go    # Exception handling
│   ├── frame.go        # Stack frames, operand stack, local vars
│   ├── heap.go         # Heap and GC
│   ├── jvm.go          # JVM state and monitors
│   ├── native.go       # Native method registry
│   ├── object.go       # Object representation
│   └── thread.go       # Thread management
├── interpreter/        # Bytecode interpreter
│   ├── interpreter.go  # Main interpreter loop
│   └── opcodes.go      # Opcode definitions
└── main.go             # Entry point
```

## Limitations

This is an educational implementation with some limitations:
- Single-threaded execution (but synchronized blocks work)
- No class loading from JAR files
- No full inheritance hierarchy resolution
- Limited native method coverage
- Simplified garbage collection

## Test Examples

The `examples/` directory contains test files:
- `HelloWorld.java` - Basic output
- `Calculator.java` - Arithmetic and recursion
- `ArrayTest.java` - Array operations
- `ObjectTest.java` - Object creation and methods
- `ExceptionTest.java` - Exception handling
- `TypeTest.java` - instanceof and checkcast
- `NativeTest.java` - Native method calls
- `SyncTest.java` - Synchronized blocks
