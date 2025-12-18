public class ArrayTest {
    public static void main(String[] args) {
        // Test primitive int array
        System.out.println("=== Int Array Test ===");
        int[] nums = new int[5];
        nums[0] = 10;
        nums[1] = 20;
        nums[2] = 30;
        nums[3] = 40;
        nums[4] = 50;
        
        System.out.println(nums[0]);
        System.out.println(nums[2]);
        System.out.println(nums[4]);
        System.out.println(nums.length);
        
        // Test sum
        int sum = 0;
        for (int i = 0; i < nums.length; i++) {
            sum = sum + nums[i];
        }
        System.out.println(sum);
        
        // Test array with loop
        System.out.println("=== Loop Test ===");
        int[] squares = new int[5];
        for (int i = 0; i < 5; i++) {
            squares[i] = i * i;
        }
        for (int i = 0; i < squares.length; i++) {
            System.out.println(squares[i]);
        }
        
        // Test long array
        System.out.println("=== Long Array Test ===");
        long[] longs = new long[3];
        longs[0] = 100000000000L;
        longs[1] = 200000000000L;
        longs[2] = 300000000000L;
        System.out.println(longs[0]);
        System.out.println(longs[1]);
        System.out.println(longs[2]);
        
        // Test String array (reference array)
        System.out.println("=== String Array Test ===");
        String[] names = new String[3];
        names[0] = "Alice";
        names[1] = "Bob";
        names[2] = "Charlie";
        System.out.println(names[0]);
        System.out.println(names[1]);
        System.out.println(names[2]);
        System.out.println(names.length);
        
        System.out.println("=== All Tests Passed ===");
    }
}


