public class ObjectTest {
    // Instance fields
    private int value;
    private String name;
    
    // Constructor
    public ObjectTest() {
        this.value = 0;
        this.name = "default";
    }
    
    // Instance methods
    public int getValue() {
        return this.value;
    }
    
    public void setValue(int v) {
        this.value = v;
    }
    
    public String getName() {
        return this.name;
    }
    
    public void setName(String n) {
        this.name = n;
    }
    
    public int add(int a, int b) {
        return a + b;
    }
    
    public int doubleValue() {
        return this.value * 2;
    }
    
    public static void main(String[] args) {
        System.out.println("=== Object Test ===");
        
        // Create an object
        ObjectTest obj = new ObjectTest();
        
        // Test getter
        System.out.println("Initial value:");
        System.out.println(obj.getValue());
        
        // Test setter
        obj.setValue(42);
        System.out.println("After setValue(42):");
        System.out.println(obj.getValue());
        
        // Test method with calculation
        System.out.println("Double value:");
        System.out.println(obj.doubleValue());
        
        // Test method with parameters
        System.out.println("add(10, 20):");
        System.out.println(obj.add(10, 20));
        
        // Test string field
        System.out.println("Initial name:");
        System.out.println(obj.getName());
        
        obj.setName("TestObject");
        System.out.println("After setName:");
        System.out.println(obj.getName());
        
        System.out.println("=== All Tests Passed ===");
    }
}


