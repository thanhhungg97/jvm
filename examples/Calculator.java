/**
 * A more complex example demonstrating:
 * - Loops (for, while)
 * - Conditionals (if/else, switch-like logic)
 * - Recursion
 * - Multiple method calls
 * - Bitwise operations
 */
public class Calculator {
    public static void main(String[] args) {
        System.out.println("=== Simple JVM Calculator ===");
        
        // Test basic arithmetic
        System.out.println("10 + 5 =");
        System.out.println(add(10, 5));
        
        System.out.println("10 - 5 =");
        System.out.println(subtract(10, 5));
        
        System.out.println("10 * 5 =");
        System.out.println(multiply(10, 5));
        
        System.out.println("10 / 5 =");
        System.out.println(divide(10, 5));
        
        System.out.println("10 % 3 =");
        System.out.println(modulo(10, 3));
        
        // Test power function (iterative)
        System.out.println("2^10 =");
        System.out.println(power(2, 10));
        
        // Test fibonacci (recursive)
        System.out.println("Fibonacci(10) =");
        System.out.println(fibonacci(10));
        
        // Test prime check
        System.out.println("Is 17 prime?");
        System.out.println(isPrime(17));
        
        System.out.println("Is 18 prime?");
        System.out.println(isPrime(18));
        
        // Test GCD (Euclidean algorithm)
        System.out.println("GCD(48, 18) =");
        System.out.println(gcd(48, 18));
        
        // Test factorial (iterative)
        System.out.println("Factorial(7) =");
        System.out.println(factorial(7));
        
        // Test sum of digits
        System.out.println("Sum of digits(12345) =");
        System.out.println(sumOfDigits(12345));
        
        // Test bitwise operations
        System.out.println("5 & 3 =");
        System.out.println(bitwiseAnd(5, 3));
        
        System.out.println("5 | 3 =");
        System.out.println(bitwiseOr(5, 3));
        
        System.out.println("5 ^ 3 =");
        System.out.println(bitwiseXor(5, 3));
        
        System.out.println("5 << 2 =");
        System.out.println(leftShift(5, 2));
        
        // Count from 1 to 5
        System.out.println("Counting 1 to 5:");
        countTo(5);
        
        System.out.println("=== Done ===");
    }
    
    // Basic arithmetic
    public static int add(int a, int b) {
        return a + b;
    }
    
    public static int subtract(int a, int b) {
        return a - b;
    }
    
    public static int multiply(int a, int b) {
        return a * b;
    }
    
    public static int divide(int a, int b) {
        return a / b;
    }
    
    public static int modulo(int a, int b) {
        return a % b;
    }
    
    // Power function (iterative)
    public static int power(int base, int exp) {
        int result = 1;
        for (int i = 0; i < exp; i++) {
            result = result * base;
        }
        return result;
    }
    
    // Fibonacci (recursive)
    public static int fibonacci(int n) {
        if (n <= 1) {
            return n;
        }
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
    
    // Prime check
    public static boolean isPrime(int n) {
        if (n <= 1) return false;
        if (n <= 3) return true;
        if (n % 2 == 0) return false;
        
        for (int i = 3; i * i <= n; i = i + 2) {
            if (n % i == 0) {
                return false;
            }
        }
        return true;
    }
    
    // GCD using Euclidean algorithm
    public static int gcd(int a, int b) {
        while (b != 0) {
            int temp = b;
            b = a % b;
            a = temp;
        }
        return a;
    }
    
    // Factorial (iterative)
    public static int factorial(int n) {
        int result = 1;
        for (int i = 2; i <= n; i++) {
            result = result * i;
        }
        return result;
    }
    
    // Sum of digits
    public static int sumOfDigits(int n) {
        int sum = 0;
        while (n > 0) {
            sum = sum + (n % 10);
            n = n / 10;
        }
        return sum;
    }
    
    // Bitwise operations
    public static int bitwiseAnd(int a, int b) {
        return a & b;
    }
    
    public static int bitwiseOr(int a, int b) {
        return a | b;
    }
    
    public static int bitwiseXor(int a, int b) {
        return a ^ b;
    }
    
    public static int leftShift(int a, int n) {
        return a << n;
    }
    
    // Loop demo
    public static void countTo(int n) {
        for (int i = 1; i <= n; i++) {
            System.out.println(i);
        }
    }
}



