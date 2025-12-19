package main

import (
	"flag"
	"fmt"
	"os"
	"simplejvm/classfile"
	"simplejvm/interpreter"
	"simplejvm/runtime"
)

func main() {
	verbose := flag.Bool("v", false, "verbose mode - print executed instructions")
	debug := flag.Bool("debug", false, "enhanced frame debugging - show locals and stack")
	trace := flag.String("trace", "", "trace calls to a method (e.g., -trace fibonacci)")
	showStats := flag.Bool("stats", false, "show heap statistics after execution")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: simplejvm [-v] [-debug] [-trace method] [-stats] <classfile>")
		fmt.Println()
		fmt.Println("A minimal JVM implementation in Go")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -v              verbose mode (print bytecode execution)")
		fmt.Println("  -debug          enhanced frame debugging (locals, stack)")
		fmt.Println("  -trace method   trace calls/returns for a method")
		fmt.Println("  -stats          show heap statistics after execution")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  simplejvm HelloWorld.class")
		fmt.Println("  simplejvm -v HelloWorld.class")
		fmt.Println("  simplejvm -debug Fib6.class")
		fmt.Println("  simplejvm -trace fibonacci Calculator.class")
		fmt.Println("  simplejvm -stats ArrayTest.class")
		os.Exit(1)
	}

	classFile := args[0]

	// Parse the class file
	cf, err := classfile.ParseFile(classFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading class file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded class: %s (Java %d)\n", cf.ClassName(), cf.MajorVersion-44)
	fmt.Println("---")

	// Create JVM instance
	jvm := runtime.NewJVM()
	defer jvm.Shutdown()

	// Create interpreter with JVM
	interp := interpreter.NewInterpreterWithJVM(*verbose, jvm)

	// Enable debug mode if requested
	if *debug {
		interp.SetDebug(true)
		fmt.Println("Debug mode enabled - showing frame state")
		fmt.Println("---")
	}

	// Enable tracing if requested
	if *trace != "" {
		interp.SetTrace(*trace)
		fmt.Printf("Tracing method: %s\n", *trace)
		fmt.Println("---")
	}

	if err := interp.Execute(cf); err != nil {
		fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("---")
	fmt.Println("Execution completed.")

	// Show heap statistics if requested
	if *showStats {
		stats := jvm.GetHeap().Stats()
		fmt.Println("---")
		fmt.Println("Heap Statistics:")
		fmt.Printf("  Allocations:  %d\n", stats.AllocCount)
		fmt.Printf("  Freed:        %d\n", stats.FreeCount)
		fmt.Printf("  Live Objects: %d\n", stats.LiveObjects)
		fmt.Printf("  Heap Size:    %d bytes\n", stats.TotalBytes)
		fmt.Printf("  GC Runs:      %d\n", stats.GCRuns)
	}
}
