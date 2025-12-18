public class ExceptionTest {
    public static void main(String[] args) {
        System.out.println("=== Exception Test ===");
        
        // Test 1: Basic try-catch
        System.out.println("Test 1: Basic try-catch");
        try {
            System.out.println("Before throw");
            throw new RuntimeException();
        } catch (RuntimeException e) {
            System.out.println("Caught RuntimeException");
        }
        System.out.println("After catch");
        
        // Test 2: Try-catch with no exception
        System.out.println("Test 2: No exception");
        try {
            System.out.println("No exception here");
        } catch (Exception e) {
            System.out.println("Should not reach here");
        }
        System.out.println("After no exception");
        
        // Test 3: Exception from method call
        System.out.println("Test 3: Exception from method");
        try {
            throwingMethod();
        } catch (Exception e) {
            System.out.println("Caught exception from method");
        }
        
        // Test 4: Arithmetic exception
        System.out.println("Test 4: Division by zero");
        try {
            int x = divideByZero(10);
            System.out.println(x); // Should not reach
        } catch (ArithmeticException e) {
            System.out.println("Caught ArithmeticException");
        }
        
        System.out.println("=== All Tests Passed ===");
    }
    
    public static void throwingMethod() {
        throw new IllegalArgumentException();
    }
    
    public static int divideByZero(int x) {
        return x / 0;
    }
}


