public class NativeTest {
    public static void main(String[] args) {
        System.out.println("=== Native Methods Test ===");
        
        // Test 1: System.currentTimeMillis
        System.out.println("Test 1: System.currentTimeMillis");
        long time1 = System.currentTimeMillis();
        System.out.println("Time1 > 0:");
        System.out.println(time1 > 0);
        
        // Small delay
        int sum = 0;
        for (int i = 0; i < 1000000; i++) {
            sum += i;
        }
        
        long time2 = System.currentTimeMillis();
        System.out.println("Time2 >= Time1:");
        System.out.println(time2 >= time1);
        
        // Test 2: System.nanoTime
        System.out.println("Test 2: System.nanoTime");
        long nano1 = System.nanoTime();
        System.out.println("Nano1 > 0:");
        System.out.println(nano1 > 0);
        
        // Test 3: System.arraycopy
        System.out.println("Test 3: System.arraycopy");
        int[] src = {1, 2, 3, 4, 5};
        int[] dest = new int[5];
        System.arraycopy(src, 0, dest, 0, 5);
        System.out.println("dest[0]:");
        System.out.println(dest[0]);
        System.out.println("dest[4]:");
        System.out.println(dest[4]);
        
        // Test 4: Math.abs
        System.out.println("Test 4: Math.abs");
        int absVal = Math.abs(-42);
        System.out.println("abs(-42):");
        System.out.println(absVal);
        
        // Test 5: Math.max/min
        System.out.println("Test 5: Math.max/min");
        System.out.println("max(10, 20):");
        System.out.println(Math.max(10, 20));
        System.out.println("min(10, 20):");
        System.out.println(Math.min(10, 20));
        
        // Test 6: Thread.sleep
        System.out.println("Test 6: Thread.sleep (100ms)");
        long before = System.currentTimeMillis();
        try {
            Thread.sleep(100);
        } catch (InterruptedException e) {
            // Ignore
        }
        long after = System.currentTimeMillis();
        System.out.println("Slept at least 50ms:");
        System.out.println((after - before) >= 50);
        
        System.out.println("=== All Tests Passed ===");
    }
}


