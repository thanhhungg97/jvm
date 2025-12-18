public class SyncTest {
    private int counter = 0;
    private static int staticCounter = 0;
    
    public synchronized void increment() {
        counter++;
    }
    
    public int getCounter() {
        return counter;
    }
    
    public static synchronized void staticIncrement() {
        staticCounter++;
    }
    
    public void syncBlock() {
        synchronized(this) {
            counter += 10;
        }
    }
    
    public static void main(String[] args) {
        System.out.println("=== Synchronized Test ===");
        
        SyncTest obj = new SyncTest();
        
        // Test 1: synchronized method
        System.out.println("Test 1: synchronized method");
        obj.increment();
        obj.increment();
        obj.increment();
        System.out.println("Counter after 3 increments:");
        System.out.println(obj.getCounter());
        
        // Test 2: synchronized block
        System.out.println("Test 2: synchronized block");
        obj.syncBlock();
        System.out.println("Counter after syncBlock (+10):");
        System.out.println(obj.getCounter());
        
        // Test 3: static synchronized method
        System.out.println("Test 3: static synchronized");
        SyncTest.staticIncrement();
        SyncTest.staticIncrement();
        System.out.println("Static counter:");
        System.out.println(staticCounter);
        
        // Test 4: nested synchronized
        System.out.println("Test 4: nested synchronized");
        synchronized(obj) {
            synchronized(obj) {
                obj.counter += 5;
            }
        }
        System.out.println("Counter after nested sync (+5):");
        System.out.println(obj.getCounter());
        
        System.out.println("=== All Tests Passed ===");
    }
}

