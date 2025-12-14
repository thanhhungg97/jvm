package main

import (
	"flag"
	"fmt"
	"os"
	"simplejvm/classfile"
	"simplejvm/interpreter"
)

func main() {
	verbose := flag.Bool("v", false, "verbose mode - print executed instructions")
	trace := flag.String("trace", "", "trace calls to a method (e.g., -trace fibonacci)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: simplejvm [-v] [-trace method] <classfile>")
		fmt.Println()
		fmt.Println("A minimal JVM implementation in Go")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -v              verbose mode (print bytecode execution)")
		fmt.Println("  -trace method   trace calls/returns for a method")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  simplejvm HelloWorld.class")
		fmt.Println("  simplejvm -v HelloWorld.class")
		fmt.Println("  simplejvm -trace fibonacci Calculator.class")
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

	// Create interpreter and execute
	interp := interpreter.NewInterpreter(*verbose)

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
}
