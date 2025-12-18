package runtime

import (
	"fmt"
	"simplejvm/classfile"
)

// JavaException represents a Java exception being thrown
type JavaException struct {
	Object    *Object // The exception object
	ClassName string  // Class name for quick lookup
	Message   string  // Exception message
}

// NewJavaException creates a new exception
func NewJavaException(obj *Object, message string) *JavaException {
	className := ""
	if obj != nil && obj.Class != nil {
		className = obj.Class.ClassName()
	}
	return &JavaException{
		Object:    obj,
		ClassName: className,
		Message:   message,
	}
}

// String returns a string representation
func (e *JavaException) String() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.ClassName, e.Message)
	}
	return e.ClassName
}

// FindExceptionHandler finds an exception handler in the code attribute
// Returns the handler PC if found, or -1 if not found
func FindExceptionHandler(code *classfile.CodeAttribute, cp classfile.ConstantPool, pc int, exceptionClass string) int {
	for _, entry := range code.ExceptionTable {
		// Check if PC is in the range [startPC, endPC)
		if pc >= int(entry.StartPC) && pc < int(entry.EndPC) {
			// Check if this handler catches this exception type
			if entry.CatchType == 0 {
				// CatchType 0 means catch all (finally block)
				return int(entry.HandlerPC)
			}

			// Get the exception class name this handler catches
			catchClassName := cp.GetClassName(entry.CatchType)

			// Check if exception matches or is a subclass
			// For now, we do simple name matching
			// TODO: Handle inheritance hierarchy
			if matchesException(exceptionClass, catchClassName) {
				return int(entry.HandlerPC)
			}
		}
	}
	return -1
}

// matchesException checks if thrownClass matches or is a subclass of catchClass
func matchesException(thrownClass, catchClass string) bool {
	// Simple matching for now
	if thrownClass == catchClass {
		return true
	}

	// Handle common Java exceptions
	// java/lang/Exception catches most exceptions
	if catchClass == "java/lang/Exception" {
		return isException(thrownClass)
	}

	// java/lang/Throwable catches everything
	if catchClass == "java/lang/Throwable" {
		return true
	}

	// java/lang/RuntimeException
	if catchClass == "java/lang/RuntimeException" {
		return isRuntimeException(thrownClass)
	}

	return false
}

// isException returns true if the class is a subclass of Exception
func isException(className string) bool {
	// Common exception classes
	exceptions := map[string]bool{
		"java/lang/Exception":                      true,
		"java/lang/RuntimeException":               true,
		"java/lang/NullPointerException":           true,
		"java/lang/ArrayIndexOutOfBoundsException": true,
		"java/lang/ArithmeticException":            true,
		"java/lang/IllegalArgumentException":       true,
		"java/lang/IllegalStateException":          true,
		"java/lang/IndexOutOfBoundsException":      true,
		"java/lang/ClassCastException":             true,
		"java/lang/NumberFormatException":          true,
		"java/io/IOException":                      true,
		"java/io/FileNotFoundException":            true,
	}
	return exceptions[className]
}

// isRuntimeException returns true if the class is a RuntimeException
func isRuntimeException(className string) bool {
	runtimeExceptions := map[string]bool{
		"java/lang/RuntimeException":               true,
		"java/lang/NullPointerException":           true,
		"java/lang/ArrayIndexOutOfBoundsException": true,
		"java/lang/ArithmeticException":            true,
		"java/lang/IllegalArgumentException":       true,
		"java/lang/IllegalStateException":          true,
		"java/lang/IndexOutOfBoundsException":      true,
		"java/lang/ClassCastException":             true,
		"java/lang/NumberFormatException":          true,
	}
	return runtimeExceptions[className]
}
