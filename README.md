# SimpleJVM

A minimal JVM implementation in Go for educational purposes.

## Features

- Class file parsing (constant pool, methods, attributes)
- Bytecode interpreter with basic instructions
- Operand stack and local variables
- Simple method invocation (static methods)
- Console output via `System.out.println`

## Supported Instructions

- Constants: `iconst_*`, `bipush`, `sipush`, `ldc`
- Loads: `iload`, `iload_*`
- Stores: `istore`, `istore_*`
- Arithmetic: `iadd`, `isub`, `imul`, `idiv`, `irem`, `ineg`
- Comparisons: `if_icmp*`, `ifeq`, `ifne`, `iflt`, `ifge`, `ifgt`, `ifle`
- Control: `goto`, `return`, `ireturn`
- Invocation: `invokestatic`, `invokevirtual` (limited)

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
├── classfile/      # Class file parser
│   ├── classfile.go
│   ├── constant_pool.go
│   └── reader.go
├── runtime/        # Runtime structures
│   ├── frame.go
│   └── thread.go
├── interpreter/    # Bytecode interpreter
│   └── interpreter.go
└── main.go         # Entry point
```

## Limitations

This is an educational implementation with many limitations:
- No garbage collection
- No multithreading
- Limited native method support
- No exception handling
- Only static method calls fully supported

# jvm
