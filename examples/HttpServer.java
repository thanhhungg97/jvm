/**
 * Simple HTTP Server using SimpleJVM native methods.
 * 
 * This demonstrates a real HTTP server running in the JVM,
 * powered by Go's net/http under the hood.
 * 
 * Run with: ./simplejvm examples/HttpServer.class
 * Then visit: http://localhost:8080/
 */
public class HttpServer {
    
    // Native method declarations (implemented in Go)
    private static native boolean startServer(int port);
    private static native void stopServer();
    private static native void addRoute(int method, int path, String response, int statusCode);
    private static native boolean isRunning();
    
    // Method constants
    public static final int GET = 1;
    public static final int POST = 2;
    public static final int PUT = 3;
    public static final int DELETE = 4;
    public static final int ANY = 0;
    
    // Path constants
    public static final int PATH_ROOT = 1;
    public static final int PATH_HELLO = 2;
    public static final int PATH_DATA = 3;
    public static final int PATH_USERS = 4;
    public static final int PATH_HEALTH = 5;
    
    public static void main(String[] args) {
        System.out.println("=== SimpleJVM HTTP Server ===");
        System.out.println("");
        
        // Register routes before starting
        System.out.println("Registering routes...");
        
        // GET / - Homepage
        addRoute(GET, PATH_ROOT, "Welcome to SimpleJVM HTTP Server!", 200);
        System.out.println("  GET /");
        
        // GET /api/hello - Hello endpoint
        addRoute(GET, PATH_HELLO, "{\"message\": \"Hello from SimpleJVM!\"}", 200);
        System.out.println("  GET /api/hello");
        
        // GET /api/data - Data endpoint
        addRoute(GET, PATH_DATA, "{\"items\": [1, 2, 3], \"count\": 3}", 200);
        System.out.println("  GET /api/data");
        
        // POST /api/data - Create data
        addRoute(POST, PATH_DATA, "{\"status\": \"created\"}", 201);
        System.out.println("  POST /api/data");
        
        // GET /api/users - Users list
        addRoute(GET, PATH_USERS, "{\"users\": [\"alice\", \"bob\"]}", 200);
        System.out.println("  GET /api/users");
        
        // GET /health - Health check
        addRoute(GET, PATH_HEALTH, "{\"status\": \"healthy\"}", 200);
        System.out.println("  GET /health");
        
        System.out.println("");
        
        // Start the server
        int port = 8080;
        System.out.println("Starting server on port:");
        System.out.println(port);
        
        boolean started = startServer(port);
        
        if (started) {
            System.out.println("Server started successfully!");
            System.out.println("");
            System.out.println("Try these URLs:");
            System.out.println("  http://localhost:8080/");
            System.out.println("  http://localhost:8080/api/hello");
            System.out.println("  http://localhost:8080/api/data");
            System.out.println("  http://localhost:8080/health");
            System.out.println("");
            System.out.println("Press Ctrl+C to stop the server");
            
            // Keep running for a while (in real app, would wait for signal)
            // The server runs for 60 seconds then shuts down
            sleep(60000);
            
            stopServer();
            System.out.println("Server stopped.");
        } else {
            System.out.println("Failed to start server!");
        }
    }
    
    // Simple sleep using native Thread.sleep
    private static void sleep(long millis) {
        try {
            Thread.sleep(millis);
        } catch (Exception e) {
            // Ignore
        }
    }
}

