/**
 * Simple Fibonacci for debugging
 * Uses small input to make trace readable
 */
public class FibonacciDebug {
    public static void main(String[] args) {
        System.out.println("Computing fibonacci(5)...");
        int result = fibonacci(5);
        System.out.println("Result:");
        System.out.println(result);
    }
    
    public static int fibonacci(int n) {
        if (n <= 1) {
            return n;
        }
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
}

