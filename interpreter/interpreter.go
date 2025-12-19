// Package interpreter executes JVM bytecode.
//
// This package is organized into several files:
//   - interpreter.go: Core interpreter loop and public API
//   - opcodes.go: Opcode constant definitions
//   - instructions_const.go: Constant-pushing instructions (iconst, ldc, etc.)
//   - instructions_load.go: Load instructions (iload, aload, etc.)
//   - instructions_store.go: Store instructions (istore, astore, etc.)
//   - instructions_array.go: Array instructions (newarray, iaload, iastore, etc.)
//   - instructions_math.go: Arithmetic and bitwise instructions
//   - instructions_control.go: Control flow (branches, returns)
//   - instructions_object.go: Object operations (new, getfield, putfield, etc.)
//   - instructions_invoke.go: Method invocation instructions
//   - helpers.go: Utility functions (tracing, constants, println, etc.)
package interpreter

import (
	"fmt"
	"simplejvm/classfile"
	"simplejvm/runtime"
)

// Interpreter executes bytecode
type Interpreter struct {
	thread       *runtime.Thread
	verbose      bool
	debug        bool // Enhanced frame debugging
	trace        bool
	traceMethod  string
	staticFields map[string]any
}

// NewInterpreter creates a new interpreter (standalone mode)
func NewInterpreter(verbose bool) *Interpreter {
	return &Interpreter{
		thread:       runtime.NewThread(),
		verbose:      verbose,
		trace:        false,
		staticFields: make(map[string]any),
	}
}

// NewInterpreterWithJVM creates a new interpreter with a JVM instance
func NewInterpreterWithJVM(verbose bool, jvm *runtime.JVM) *Interpreter {
	thread := jvm.CreateThread()
	interp := &Interpreter{
		thread:       thread,
		verbose:      verbose,
		trace:        false,
		staticFields: make(map[string]any),
	}
	// Register this interpreter as the callback executor
	runtime.SetCallbackExecutor(interp)
	return interp
}

// SetTrace enables simple call/return tracing for a specific method
func (i *Interpreter) SetTrace(methodName string) {
	i.trace = true
	i.traceMethod = methodName
}

// SetDebug enables enhanced frame debugging
func (i *Interpreter) SetDebug(enabled bool) {
	i.debug = enabled
}

// Execute runs the main method of a class
func (i *Interpreter) Execute(cf *classfile.ClassFile) error {
	className := cf.ClassName()
	i.thread.LoadClass(className, cf)

	// Print constant pool in debug mode
	if i.debug {
		PrintConstantPool(cf.ConstantPool)
	}

	mainMethod := cf.GetMethod("main", "([Ljava/lang/String;)V")
	if mainMethod == nil {
		return fmt.Errorf("main method not found in class %s", className)
	}

	frame := runtime.NewFrame(i.thread, mainMethod, cf)
	if frame == nil {
		return fmt.Errorf("could not create frame for main method")
	}

	frame.LocalVars.SetRef(0, nil) // args placeholder
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
			i.thread.PopFrame()
			continue
		}

		pc := frame.PC
		opcode := frame.ReadU1()
		methodName := frame.Method.Name(frame.Class.ConstantPool)

		if i.debug {
			i.printFrameDebug(frame, pc, opcode)
		} else if i.verbose {
			fmt.Printf("[%s] PC=%d opcode=0x%02X\n", methodName, pc, opcode)
		}

		if err := i.executeInstruction(frame, opcode); err != nil {
			return err
		}
	}
	return nil
}

// InvokeMethod invokes a method on an object (implements CallbackExecutor)
func (i *Interpreter) InvokeMethod(obj interface{}, methodName, descriptor string) error {
	// Get the object and its class
	runtimeObj, ok := obj.(*runtime.Object)
	if !ok {
		return fmt.Errorf("InvokeMethod: expected *runtime.Object, got %T", obj)
	}

	if runtimeObj.Class == nil {
		return fmt.Errorf("InvokeMethod: object has no class")
	}

	// Find the method
	method := runtimeObj.Class.GetMethod(methodName, descriptor)
	if method == nil {
		return fmt.Errorf("InvokeMethod: method %s%s not found in %s",
			methodName, descriptor, runtimeObj.Class.ClassName())
	}

	// Create a new frame for the method
	frame := runtime.NewFrame(i.thread, method, runtimeObj.Class)
	if frame == nil {
		return fmt.Errorf("InvokeMethod: could not create frame for %s", methodName)
	}

	// Set 'this' as local var 0
	frame.LocalVars.SetRef(0, obj)

	// Push the frame and execute
	i.thread.PushFrame(frame)

	// Run until this frame completes
	return i.runUntilFrameCompletes(frame)
}

// InvokeRunnable invokes the run() method on a Runnable object (implements CallbackExecutor)
func (i *Interpreter) InvokeRunnable(runnable interface{}) error {
	return i.InvokeMethod(runnable, "run", "()V")
}

// runUntilFrameCompletes runs the interpreter until a specific frame completes
func (i *Interpreter) runUntilFrameCompletes(targetFrame *runtime.Frame) error {
	for !i.thread.IsStackEmpty() {
		frame := i.thread.CurrentFrame()
		if frame == nil {
			break
		}

		// Check if we've returned from the target frame
		if frame != targetFrame && !i.thread.ContainsFrame(targetFrame) {
			// Target frame has been popped, we're done
			break
		}

		if frame.PC >= len(frame.Code) {
			i.thread.PopFrame()
			continue
		}

		pc := frame.PC
		opcode := frame.ReadU1()

		if i.debug {
			i.printFrameDebug(frame, pc, opcode)
		} else if i.verbose {
			methodName := frame.Method.Name(frame.Class.ConstantPool)
			fmt.Printf("[%s] PC=%d opcode=0x%02X\n", methodName, pc, opcode)
		}

		if err := i.executeInstruction(frame, opcode); err != nil {
			return err
		}
	}
	return nil
}

// executeInstruction dispatches to the appropriate instruction handler based on opcode category
func (i *Interpreter) executeInstruction(frame *runtime.Frame, opcode uint8) error {
	switch Category(opcode) {
	case CategoryConst:
		i.executeConstInstruction(frame, opcode)
		return nil

	case CategoryLoad:
		i.executeLoadInstruction(frame, opcode)
		return nil

	case CategoryStore:
		i.executeStoreInstruction(frame, opcode)
		return nil

	case CategoryMath:
		_, err := i.executeMathInstruction(frame, opcode)
		return err

	case CategoryControl:
		i.executeControlInstruction(frame, opcode)
		return nil

	case CategoryArray:
		_, err := i.executeArrayInstruction(frame, opcode)
		return err

	case CategoryObject:
		_, err := i.executeObjectInstruction(frame, opcode)
		return err

	case CategoryInvoke:
		_, err := i.executeInvokeInstruction(frame, opcode)
		return err

	default:
		return fmt.Errorf("unimplemented opcode: 0x%02X at PC=%d", opcode, frame.PC-1)
	}
}
