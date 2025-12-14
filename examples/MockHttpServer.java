/**
 * Mock HTTP Server - demonstrates what a real server would look like
 * Uses only integers since String parameters aren't fully supported yet
 */
public class MockHttpServer {
    
    public static void main(String[] args) {
        System.out.println("=== Mock HTTP Server ===");
        System.out.println("Starting server on port 8080");
        System.out.println("Listening for connections...");
        
        // Simulate handling requests
        // Method: 1=GET, 2=POST, 3=PUT, 4=DELETE
        // Path: 1=/, 2=/api/users, 3=/api/data
        
        handleRequest(1, 1, 1);  // Request 1: GET /
        handleRequest(2, 1, 2);  // Request 2: GET /api/users
        handleRequest(3, 2, 3);  // Request 3: POST /api/data
        handleRequest(4, 4, 2);  // Request 4: DELETE /api/users
        
        System.out.println("Server shutting down");
    }
    
    public static void handleRequest(int requestId, int method, int path) {
        System.out.println("--- Request ---");
        System.out.println("Request ID:");
        System.out.println(requestId);
        
        printMethod(method);
        printPath(path);
        
        // Route and get status code
        int statusCode = route(method, path);
        
        System.out.println("Response:");
        System.out.println(statusCode);
    }
    
    public static void printMethod(int method) {
        System.out.println("Method:");
        // 1=GET, 2=POST, 3=PUT, 4=DELETE
        if (method == 1) {
            System.out.println("GET");
        } else if (method == 2) {
            System.out.println("POST");
        } else if (method == 3) {
            System.out.println("PUT");
        } else if (method == 4) {
            System.out.println("DELETE");
        }
    }
    
    public static void printPath(int path) {
        System.out.println("Path:");
        if (path == 1) {
            System.out.println("/");
        } else if (path == 2) {
            System.out.println("/api/users");
        } else if (path == 3) {
            System.out.println("/api/data");
        }
    }
    
    public static int route(int method, int path) {
        // GET requests
        if (method == 1) {
            if (path == 1) {
                return 200;  // OK - homepage
            } else if (path == 2) {
                return 200;  // OK - users list
            }
            return 404;  // Not Found
        }
        
        // POST requests
        if (method == 2) {
            if (path == 3) {
                return 201;  // Created
            }
            return 400;  // Bad Request
        }
        
        // DELETE requests
        if (method == 4) {
            if (path == 2) {
                return 204;  // No Content (deleted)
            }
            return 404;  // Not Found
        }
        
        return 405;  // Method Not Allowed
    }
}
