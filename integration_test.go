package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"simplejvm/classfile"
	"simplejvm/interpreter"
	"simplejvm/runtime"
)

// TestJavaExamples runs all Java example files and verifies they execute successfully
func TestJavaExamples(t *testing.T) {
	examples := []struct {
		name           string
		classFile      string
		expectedOutput []string // Substrings that should appear in output
	}{
		{
			name:      "HelloWorld",
			classFile: "examples/HelloWorld.class",
			expectedOutput: []string{
				"Hello from SimpleJVM!",
				"30",
				"120",
			},
		},
		{
			name:      "Calculator",
			classFile: "examples/Calculator.class",
			expectedOutput: []string{
				"=== Simple JVM Calculator ===",
				"10 + 5 =",
				"15",
				"Fibonacci(10) =",
				"55",
				"=== Done ===",
			},
		},
		{
			name:      "ArrayTest",
			classFile: "examples/ArrayTest.class",
			expectedOutput: []string{
				"=== Int Array Test ===",
				"=== Loop Test ===",
				"=== Long Array Test ===",
				"=== String Array Test ===",
				"=== All Tests Passed ===",
			},
		},
		{
			name:      "ObjectTest",
			classFile: "examples/ObjectTest.class",
			expectedOutput: []string{
				"=== Object Test ===",
				"Initial value:",
				"42",
				"TestObject",
				"=== All Tests Passed ===",
			},
		},
		{
			name:      "ExceptionTest",
			classFile: "examples/ExceptionTest.class",
			expectedOutput: []string{
				"=== Exception Test ===",
				"Caught RuntimeException",
				"Caught exception from method",
				"Caught ArithmeticException",
				"=== All Tests Passed ===",
			},
		},
		{
			name:      "TypeTest",
			classFile: "examples/TypeTest.class",
			expectedOutput: []string{
				"=== Type Test ===",
				"t is TypeTest",
				"o is Object",
				"null is not instanceof anything",
				"Cast succeeded",
				"=== All Tests Passed ===",
			},
		},
		{
			name:      "NativeTest",
			classFile: "examples/NativeTest.class",
			expectedOutput: []string{
				"=== Native Methods Test ===",
				"Time1 > 0:",
				"true",
				"dest[0]:",
				"1",
				"abs(-42):",
				"42",
				"=== All Tests Passed ===",
			},
		},
		{
			name:      "SyncTest",
			classFile: "examples/SyncTest.class",
			expectedOutput: []string{
				"=== Synchronized Test ===",
				"Counter after 3 increments:",
				"3",
				"Static counter:",
				"2",
				"=== All Tests Passed ===",
			},
		},
	}

	for _, tt := range examples {
		t.Run(tt.name, func(t *testing.T) {
			// Check if class file exists
			if _, err := os.Stat(tt.classFile); os.IsNotExist(err) {
				t.Skipf("Class file %s not found, skipping", tt.classFile)
			}

			// Capture stdout
			output := captureOutput(t, func() error {
				return runClass(tt.classFile)
			})

			// Check for expected output
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Output missing expected string: %q\nOutput:\n%s", expected, output)
				}
			}
		})
	}
}

// runClass executes a class file and returns any error
func runClass(classFile string) error {
	cf, err := classfile.ParseFile(classFile)
	if err != nil {
		return fmt.Errorf("failed to parse class file: %w", err)
	}

	jvm := runtime.NewJVM()
	defer jvm.Shutdown()

	interp := interpreter.NewInterpreterWithJVM(false, jvm)
	return interp.Execute(cf)
}

// captureOutput captures stdout during function execution
func captureOutput(t *testing.T, fn func() error) string {
	// Save original stdout
	oldStdout := os.Stdout

	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Replace stdout
	os.Stdout = w

	// Run the function
	fnErr := fn()

	// Close writer and restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	r.Close()

	if fnErr != nil {
		t.Errorf("Execution error: %v", fnErr)
	}

	return buf.String()
}

// TestClassFileParsing tests that all example class files can be parsed
func TestClassFileParsing(t *testing.T) {
	classFiles, err := filepath.Glob("examples/*.class")
	if err != nil {
		t.Fatalf("Failed to glob class files: %v", err)
	}

	if len(classFiles) == 0 {
		t.Skip("No class files found in examples/")
	}

	for _, classFile := range classFiles {
		name := filepath.Base(classFile)
		t.Run(name, func(t *testing.T) {
			cf, err := classfile.ParseFile(classFile)
			if err != nil {
				t.Errorf("Failed to parse %s: %v", classFile, err)
				return
			}

			// Verify basic class file structure
			if cf.Magic != 0xCAFEBABE {
				t.Errorf("Invalid magic number: 0x%X", cf.Magic)
			}

			className := cf.ClassName()
			if className == "" {
				t.Error("Class name is empty")
			}

			// Check that main method exists (for test classes)
			mainMethod := cf.GetMethod("main", "([Ljava/lang/String;)V")
			if mainMethod == nil {
				t.Logf("Note: %s has no main method", name)
			}
		})
	}
}

// TestCompileAndRun compiles and runs a simple Java program
func TestCompileAndRun(t *testing.T) {
	// Skip if javac is not available
	if _, err := exec.LookPath("javac"); err != nil {
		t.Skip("javac not found, skipping compile test")
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "simplejvm-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write a simple Java file
	javaCode := `
public class SimpleTest {
    public static void main(String[] args) {
        System.out.println("SimpleTest OK");
        System.out.println(add(3, 4));
    }
    public static int add(int a, int b) {
        return a + b;
    }
}
`
	javaFile := filepath.Join(tmpDir, "SimpleTest.java")
	if err := os.WriteFile(javaFile, []byte(javaCode), 0644); err != nil {
		t.Fatalf("Failed to write Java file: %v", err)
	}

	// Compile
	cmd := exec.Command("javac", javaFile)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to compile: %v\n%s", err, output)
	}

	// Run with our JVM
	classFile := filepath.Join(tmpDir, "SimpleTest.class")
	output := captureOutput(t, func() error {
		return runClass(classFile)
	})

	if !strings.Contains(output, "SimpleTest OK") {
		t.Errorf("Missing expected output 'SimpleTest OK'\nGot: %s", output)
	}
	if !strings.Contains(output, "7") {
		t.Errorf("Missing expected output '7'\nGot: %s", output)
	}
}

// BenchmarkExecution benchmarks the execution of Calculator.class
func BenchmarkCalculator(b *testing.B) {
	classFile := "examples/Calculator.class"
	if _, err := os.Stat(classFile); os.IsNotExist(err) {
		b.Skip("Calculator.class not found")
	}

	cf, err := classfile.ParseFile(classFile)
	if err != nil {
		b.Fatalf("Failed to parse: %v", err)
	}

	// Suppress output during benchmark
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jvm := runtime.NewJVM()
		interp := interpreter.NewInterpreterWithJVM(false, jvm)
		interp.Execute(cf)
		jvm.Shutdown()
	}
}
