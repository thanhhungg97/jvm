public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello from SimpleJVM!");
        
        // Test arithmetic
        int a = 10;
        int b = 20;
        int sum = add(a, b);
        System.out.println(sum);
        
        // Test multiplication
        int product = multiply(5, 7);
        System.out.println(product);
        
        // Test conditionals
        int max = max(42, 17);
        System.out.println(max);
        
        // Test loop
        int factorial = factorial(5);
        System.out.println(factorial);
    }
    
    public static int add(int x, int y) {
        return x + y;
    }
    
    public static int multiply(int x, int y) {
        return x * y;
    }
    
    public static int max(int x, int y) {
        if (x > y) {
            return x;
        } else {
            return y;
        }
    }
    
    public static int factorial(int n) {
        int result = 1;
        for (int i = 2; i <= n; i++) {
            result = result * i;
        }
        return result;
    }
}

