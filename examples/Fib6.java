public class Fib6 {
    public static void main(String[] args) {
        System.out.println("fibonacci(6) =");
        int result = fibonacci(6);
        System.out.println(result);
    }
    
    public static int fibonacci(int n) {
        if (n <= 1) {
            return n;
        }
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
}

