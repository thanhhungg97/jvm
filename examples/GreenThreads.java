/**
 * Green Threads (Fibers) demonstration in SimpleJVM.
 * 
 * This shows lightweight user-space threads that are scheduled
 * cooperatively, similar to Go's goroutines or Erlang's processes.
 * 
 * Run with: ./simplejvm examples/GreenThreads.class
 */
public class GreenThreads {
    
    // Native fiber operations (implemented in Go)
    private static native long spawn(int taskId, String name);
    private static native void yield();
    private static native void sleep(long millis);
    private static native void join(long fiberId);
    private static native boolean isAlive(long fiberId);
    private static native long current();
    private static native int count();
    private static native void printStats();
    
    public static void main(String[] args) {
        System.out.println("=== Green Threads Demo ===");
        System.out.println("");
        
        // Test 1: Spawn multiple fibers
        System.out.println("Test 1: Spawning 3 concurrent fibers");
        System.out.println("---");
        
        long fiber1 = spawn(2, "alpha");
        long fiber2 = spawn(3, "beta");
        long fiber3 = spawn(1, "gamma");
        
        System.out.println("Spawned fiber IDs:");
        System.out.println(fiber1);
        System.out.println(fiber2);
        System.out.println(fiber3);
        System.out.println("");
        
        // Check active count
        System.out.println("Active fibers:");
        System.out.println(count());
        System.out.println("");
        
        // Wait for all fibers to complete
        System.out.println("Waiting for fibers to complete...");
        join(fiber1);
        join(fiber2);
        join(fiber3);
        
        System.out.println("");
        System.out.println("All fibers completed!");
        System.out.println("");
        
        // Print statistics
        printStats();
        System.out.println("");
        
        // Test 2: Check isAlive
        System.out.println("Test 2: Fiber lifecycle");
        System.out.println("---");
        
        long shortFiber = spawn(1, "short");
        System.out.println("Fiber alive before join:");
        System.out.println(isAlive(shortFiber) ? 1 : 0);
        
        join(shortFiber);
        
        System.out.println("Fiber alive after join:");
        System.out.println(isAlive(shortFiber) ? 1 : 0);
        System.out.println("");
        
        // Test 3: Yielding
        System.out.println("Test 3: Cooperative yielding");
        System.out.println("---");
        
        long yieldFiber1 = spawn(2, "yield-a");
        long yieldFiber2 = spawn(2, "yield-b");
        
        // Main thread yields a few times
        for (int i = 0; i < 3; i++) {
            System.out.println("Main yielding...");
            yield();
        }
        
        join(yieldFiber1);
        join(yieldFiber2);
        
        System.out.println("");
        System.out.println("=== Demo Complete ===");
        printStats();
    }
}

