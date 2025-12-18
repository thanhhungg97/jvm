public class TypeTest {
    public static void main(String[] args) {
        System.out.println("=== Type Test ===");
        
        // Test 1: instanceof with same type
        System.out.println("Test 1: instanceof same type");
        TypeTest t = new TypeTest();
        if (t instanceof TypeTest) {
            System.out.println("t is TypeTest");
        }
        
        // Test 2: instanceof with Object
        System.out.println("Test 2: instanceof Object");
        Object o = new TypeTest();
        if (o instanceof Object) {
            System.out.println("o is Object");
        }
        
        // Test 3: instanceof with null
        System.out.println("Test 3: instanceof null");
        TypeTest nullRef = null;
        if (nullRef instanceof TypeTest) {
            System.out.println("Should not print");
        } else {
            System.out.println("null is not instanceof anything");
        }
        
        // Test 4: checkcast (implicit in assignment)
        System.out.println("Test 4: checkcast");
        Object obj = new TypeTest();
        TypeTest casted = (TypeTest) obj;
        System.out.println("Cast succeeded");
        
        // Test 5: Array instanceof
        System.out.println("Test 5: Array instanceof");
        int[] arr = new int[5];
        if (arr instanceof Object) {
            System.out.println("array is Object");
        }
        
        System.out.println("=== All Tests Passed ===");
    }
}


