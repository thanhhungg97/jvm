/**
 * Event Loop Demo - Node.js-style event loop in SimpleJVM.
 * 
 * Demonstrates cooperative multitasking using an event loop pattern:
 * - Task queue (FIFO)
 * - Timers (setTimeout)
 * - TRUE LAMBDAS - pass actual Runnable objects!
 * 
 * Run with: ./simplejvm examples/EventLoopDemo.class
 */
public class EventLoopDemo {
    
    // Native event loop operations (implemented in Go)
    private static native void submit(int taskId, String name);
    private static native void submitRunnable(Runnable r);
    private static native void setTimeout(int taskId, String name, long delayMs);
    private static native void setTimeoutRunnable(Runnable r, long delayMs);
    private static native void run();
    private static native void stop();
    private static native boolean isRunning();
    private static native void printStats();
    private static native void reset();
    
    public static void main(String[] args) {
        System.out.println("=== Event Loop Demo with True Lambdas ===");
        System.out.println("");
        
        // Test 1: Basic ID-based tasks (old way)
        System.out.println("Test 1: Basic Tasks (old way - ID based)");
        System.out.println("---");
        reset();
        
        submit(1, "task-1");
        submit(2, "task-2");
        
        run();
        System.out.println("");
        
        // Test 2: TRUE LAMBDAS with Runnable!
        System.out.println("Test 2: True Lambdas (Runnable objects)");
        System.out.println("---");
        reset();
        
        // Submit Runnable tasks
        submitRunnable(new PrintTask("Lambda task 1 executed!"));
        submitRunnable(new PrintTask("Lambda task 2 executed!"));
        submitRunnable(new PrintTask("Lambda task 3 executed!"));
        
        System.out.println("Submitted 3 lambda tasks");
        run();
        System.out.println("");
        
        // Test 3: Delayed lambdas with setTimeout
        System.out.println("Test 3: Delayed Lambdas (setTimeout)");
        System.out.println("---");
        reset();
        
        submitRunnable(new PrintTask("Immediate task A"));
        setTimeoutRunnable(new PrintTask("Delayed 100ms"), 100);
        setTimeoutRunnable(new PrintTask("Delayed 50ms"), 50);
        submitRunnable(new PrintTask("Immediate task B"));
        
        System.out.println("Submitted immediate + delayed tasks");
        run();
        System.out.println("");
        
        // Test 4: Computation in lambda
        System.out.println("Test 4: Computation in Lambda");
        System.out.println("---");
        reset();
        
        submitRunnable(new FactorialTask(5));
        submitRunnable(new FactorialTask(7));
        submitRunnable(new FactorialTask(10));
        
        run();
        System.out.println("");
        
        System.out.println("=== Demo Complete ===");
        printStats();
    }
}

// A simple task that prints a message
class PrintTask implements Runnable {
    private String message;
    
    public PrintTask(String msg) {
        this.message = msg;
    }
    
    public void run() {
        System.out.println(message);
    }
}

// A task that computes factorial and prints result
class FactorialTask implements Runnable {
    private int n;
    
    public FactorialTask(int n) {
        this.n = n;
    }
    
    public void run() {
        int result = factorial(n);
        // Avoid string concatenation, just print the values
        System.out.print("factorial(");
        System.out.print(n);
        System.out.print(") = ");
        System.out.println(result);
    }
    
    private int factorial(int x) {
        if (x <= 1) return 1;
        return x * factorial(x - 1);
    }
}
