.PHONY: build run clean test

# Use local cache to avoid sandbox permission issues
export GOCACHE := $(PWD)/.cache

build:
	go build -o simplejvm

run: build
	./simplejvm examples/HelloWorld.class

run-verbose: build
	./simplejvm -v examples/HelloWorld.class

test: build
	cd examples && javac HelloWorld.java
	./simplejvm examples/HelloWorld.class

clean:
	rm -f simplejvm
	rm -f examples/*.class
	rm -rf .cache

