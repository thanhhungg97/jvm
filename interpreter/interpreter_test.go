package interpreter

import (
	"simplejvm/classfile"
	"simplejvm/runtime"
	"testing"
)

// ==================== Opcode Tests ====================

func TestOpcodeConstants(t *testing.T) {
	tests := []struct {
		name   string
		opcode uint8
		want   uint8
	}{
		{"NOP", NOP, 0x00},
		{"ICONST_0", ICONST_0, 0x03},
		{"ICONST_5", ICONST_5, 0x08},
		{"BIPUSH", BIPUSH, 0x10},
		{"ILOAD", ILOAD, 0x15},
		{"ISTORE", ISTORE, 0x36},
		{"IADD", IADD, 0x60},
		{"ISUB", ISUB, 0x64},
		{"IMUL", IMUL, 0x68},
		{"IDIV", IDIV, 0x6C},
		{"GOTO", GOTO, 0xA7},
		{"IRETURN", IRETURN, 0xAC},
		{"RETURN", RETURN, 0xB1},
		{"INVOKESTATIC", INVOKESTATIC, 0xB8},
		{"INVOKEVIRTUAL", INVOKEVIRTUAL, 0xB6},
		{"NEW", NEW, 0xBB},
		{"NEWARRAY", NEWARRAY, 0xBC},
		{"ATHROW", ATHROW, 0xBF},
		{"MONITORENTER", MONITORENTER, 0xC2},
		{"MONITOREXIT", MONITOREXIT, 0xC3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.opcode != tt.want {
				t.Errorf("%s = 0x%02X, want 0x%02X", tt.name, tt.opcode, tt.want)
			}
		})
	}
}

// ==================== Helper Function Tests ====================

func TestCountArgs(t *testing.T) {
	tests := []struct {
		descriptor string
		want       int
	}{
		{"()V", 0},
		{"(I)V", 1},
		{"(II)I", 2},
		{"(IJ)V", 2},
		{"(Ljava/lang/String;)V", 1},
		{"(ILjava/lang/String;I)V", 3},
		{"([I)V", 1},
		{"([Ljava/lang/Object;)V", 1},
		{"(II[BLjava/lang/String;)I", 4},
	}

	for _, tt := range tests {
		t.Run(tt.descriptor, func(t *testing.T) {
			got := countArgs(tt.descriptor)
			if got != tt.want {
				t.Errorf("countArgs(%q) = %d, want %d", tt.descriptor, got, tt.want)
			}
		})
	}
}

func TestParseArgTypes(t *testing.T) {
	tests := []struct {
		descriptor string
		want       []byte
	}{
		{"()V", nil},
		{"(I)V", []byte{'I'}},
		{"(IJ)V", []byte{'I', 'J'}},
		{"(Ljava/lang/String;)V", []byte{'L'}},
		{"([I)V", []byte{'['}},
		{"(ILjava/lang/Object;[B)V", []byte{'I', 'L', '['}},
	}

	for _, tt := range tests {
		t.Run(tt.descriptor, func(t *testing.T) {
			got := parseArgTypes(tt.descriptor)
			if len(got) != len(tt.want) {
				t.Errorf("parseArgTypes(%q) len = %d, want %d", tt.descriptor, len(got), len(tt.want))
				return
			}
			for i, b := range got {
				if b != tt.want[i] {
					t.Errorf("parseArgTypes(%q)[%d] = %c, want %c", tt.descriptor, i, b, tt.want[i])
				}
			}
		})
	}
}

func TestExtractClassName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Object<java/lang/String>", "java/lang/String"},
		{"Object<java/lang/RuntimeException>", "java/lang/RuntimeException"},
		{"java/lang/Object", "java/lang/Object"},
		{"", ""},
		{"Object<>", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractClassName(tt.input)
			if got != tt.want {
				t.Errorf("extractClassName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ==================== Interpreter Creation Tests ====================

func TestNewInterpreter(t *testing.T) {
	interp := NewInterpreter(false)
	if interp == nil {
		t.Fatal("NewInterpreter returned nil")
		return // unreachable but satisfies staticcheck
	}
	if interp.thread == nil {
		t.Error("Interpreter thread is nil")
	}
	if interp.staticFields == nil {
		t.Error("Interpreter staticFields is nil")
	}
}

func TestNewInterpreterWithJVM(t *testing.T) {
	jvm := runtime.NewJVM()
	interp := NewInterpreterWithJVM(true, jvm)

	if interp == nil {
		t.Fatal("NewInterpreterWithJVM returned nil")
		return // unreachable but satisfies staticcheck
	}
	if interp.thread == nil {
		t.Error("Interpreter thread is nil")
	}
	if !interp.verbose {
		t.Error("Interpreter verbose should be true")
	}
}

func TestSetTrace(t *testing.T) {
	interp := NewInterpreter(false)
	interp.SetTrace("testMethod")

	if !interp.trace {
		t.Error("trace should be true")
	}
	if interp.traceMethod != "testMethod" {
		t.Errorf("traceMethod = %q, want testMethod", interp.traceMethod)
	}
}

// ==================== Exception Handling Tests ====================

func TestFindExceptionHandler(t *testing.T) {
	// Create a mock code attribute with exception table
	code := &classfile.CodeAttribute{
		ExceptionTable: []*classfile.ExceptionTableEntry{
			{StartPC: 0, EndPC: 10, HandlerPC: 20, CatchType: 0}, // catch all
		},
	}

	// Test PC in range
	handler := runtime.FindExceptionHandler(code, nil, 5, "java/lang/Exception")
	if handler != 20 {
		t.Errorf("FindExceptionHandler = %d, want 20", handler)
	}

	// Test PC out of range
	handler = runtime.FindExceptionHandler(code, nil, 15, "java/lang/Exception")
	if handler != -1 {
		t.Errorf("FindExceptionHandler = %d, want -1", handler)
	}
}
