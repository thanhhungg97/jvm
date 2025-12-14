/**
 * Classic algorithms implemented in Java
 * Tests recursion, loops, and method calls
 */
public class Algorithms {
    public static void main(String[] args) {
        System.out.println("=== Algorithm Tests ===");
        
        // Binary search simulation (without arrays)
        // Finding how many steps to find 42 in range 1-100
        System.out.println("Binary search steps for 42 in 1-100:");
        System.out.println(binarySearchSteps(1, 100, 42));
        
        // Collatz conjecture (3n+1 problem)
        System.out.println("Collatz steps for 27:");
        System.out.println(collatzSteps(27));
        
        // Count set bits (Brian Kernighan's algorithm)
        System.out.println("Set bits in 255:");
        System.out.println(countSetBits(255));
        
        System.out.println("Set bits in 128:");
        System.out.println(countSetBits(128));
        
        // Integer square root (binary search)
        System.out.println("sqrt(144) =");
        System.out.println(intSqrt(144));
        
        System.out.println("sqrt(1000000) =");
        System.out.println(intSqrt(1000000));
        
        // Ackermann function (limited - very recursive)
        System.out.println("Ackermann(2, 3) =");
        System.out.println(ackermann(2, 3));
        
        // Sum of first N numbers using formula
        System.out.println("Sum 1 to 100 =");
        System.out.println(sumFormula(100));
        
        // Sum of first N numbers using loop
        System.out.println("Sum 1 to 100 (loop) =");
        System.out.println(sumLoop(100));
        
        // Test negative numbers
        System.out.println("abs(-42) =");
        System.out.println(abs(-42));
        
        // Min/Max
        System.out.println("min(7, 3) =");
        System.out.println(min(7, 3));
        
        System.out.println("max(7, 3) =");
        System.out.println(max(7, 3));
        
        // Clamp
        System.out.println("clamp(150, 0, 100) =");
        System.out.println(clamp(150, 0, 100));
        
        System.out.println("=== Done ===");
    }
    
    // Simulates binary search - returns number of steps
    public static int binarySearchSteps(int low, int high, int target) {
        int steps = 0;
        while (low <= high) {
            int mid = low + (high - low) / 2;
            steps = steps + 1;
            if (mid == target) {
                return steps;
            } else if (mid < target) {
                low = mid + 1;
            } else {
                high = mid - 1;
            }
        }
        return steps;
    }
    
    // Collatz conjecture - count steps to reach 1
    public static int collatzSteps(int n) {
        int steps = 0;
        while (n != 1) {
            if (n % 2 == 0) {
                n = n / 2;
            } else {
                n = 3 * n + 1;
            }
            steps = steps + 1;
        }
        return steps;
    }
    
    // Count set bits using Brian Kernighan's algorithm
    public static int countSetBits(int n) {
        int count = 0;
        while (n > 0) {
            n = n & (n - 1);  // Clear lowest set bit
            count = count + 1;
        }
        return count;
    }
    
    // Integer square root using binary search
    public static int intSqrt(int n) {
        if (n < 2) return n;
        
        int low = 1;
        int high = n / 2;
        int result = 1;
        
        while (low <= high) {
            int mid = low + (high - low) / 2;
            int square = mid * mid;
            
            if (square == n) {
                return mid;
            } else if (square < n) {
                result = mid;
                low = mid + 1;
            } else {
                high = mid - 1;
            }
        }
        return result;
    }
    
    // Ackermann function (careful with large inputs!)
    public static int ackermann(int m, int n) {
        if (m == 0) {
            return n + 1;
        } else if (n == 0) {
            return ackermann(m - 1, 1);
        } else {
            return ackermann(m - 1, ackermann(m, n - 1));
        }
    }
    
    // Sum using Gauss formula: n*(n+1)/2
    public static int sumFormula(int n) {
        return n * (n + 1) / 2;
    }
    
    // Sum using loop
    public static int sumLoop(int n) {
        int sum = 0;
        for (int i = 1; i <= n; i++) {
            sum = sum + i;
        }
        return sum;
    }
    
    // Absolute value
    public static int abs(int n) {
        if (n < 0) {
            return -n;
        }
        return n;
    }
    
    // Minimum
    public static int min(int a, int b) {
        if (a < b) {
            return a;
        }
        return b;
    }
    
    // Maximum
    public static int max(int a, int b) {
        if (a > b) {
            return a;
        }
        return b;
    }
    
    // Clamp value to range
    public static int clamp(int value, int minVal, int maxVal) {
        if (value < minVal) {
            return minVal;
        }
        if (value > maxVal) {
            return maxVal;
        }
        return value;
    }
}

