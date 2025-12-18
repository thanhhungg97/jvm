package runtime

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HTTPServer represents a simple HTTP server instance
type HTTPServer struct {
	server   *http.Server
	mux      *http.ServeMux
	handlers map[string]*HTTPHandler
	running  bool
	mu       sync.RWMutex
}

// HTTPHandler represents a registered handler with its response
type HTTPHandler struct {
	Path        string
	Method      string
	Response    string
	StatusCode  int
	ContentType string
}

// Global HTTP server instance
var httpServer *HTTPServer

// NewHTTPServer creates a new HTTP server
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{
		mux:      http.NewServeMux(),
		handlers: make(map[string]*HTTPHandler),
	}
}

// Start starts the HTTP server on the given port
func (s *HTTPServer) Start(port int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server already running")
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.mux,
	}

	// Register catch-all handler that routes to our handlers
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		key := r.Method + ":" + r.URL.Path
		handler, ok := s.handlers[key]
		s.mu.RUnlock()

		if !ok {
			// Try wildcard
			s.mu.RLock()
			handler, ok = s.handlers["*:"+r.URL.Path]
			s.mu.RUnlock()
		}

		if ok {
			if handler.ContentType != "" {
				w.Header().Set("Content-Type", handler.ContentType)
			} else {
				w.Header().Set("Content-Type", "text/plain")
			}
			w.WriteHeader(handler.StatusCode)
			io.WriteString(w, handler.Response)
		} else {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "Not Found")
		}
	})

	s.running = true

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Give the server time to start
	time.Sleep(50 * time.Millisecond)
	return nil
}

// Stop stops the HTTP server
func (s *HTTPServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.running = false
	return s.server.Shutdown(ctx)
}

// RegisterHandler registers a handler for a path
func (s *HTTPServer) RegisterHandler(method, path, response string, statusCode int, contentType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := method + ":" + path
	s.handlers[key] = &HTTPHandler{
		Path:        path,
		Method:      method,
		Response:    response,
		StatusCode:  statusCode,
		ContentType: contentType,
	}
}

// =============== Native Method Implementations ===============

func init() {
	// Register HTTP server natives for the HttpServer class (no package)
	Natives.Register("HttpServer", "startServer", "(I)Z", nativeSimpleHttpStart)
	Natives.Register("HttpServer", "stopServer", "()V", nativeSimpleHttpStop)
	Natives.Register("HttpServer", "addRoute", "(IILjava/lang/String;I)V", nativeSimpleHttpAddRoute)
	Natives.Register("HttpServer", "isRunning", "()Z", nativeSimpleHttpIsRunning)

	// Also register with full package path for packaged classes
	Natives.Register("simplejvm/http/HttpServer", "startServer", "(I)Z", nativeSimpleHttpStart)
	Natives.Register("simplejvm/http/HttpServer", "stopServer", "()V", nativeSimpleHttpStop)
	Natives.Register("simplejvm/http/HttpServer", "addRoute", "(IILjava/lang/String;I)V", nativeSimpleHttpAddRoute)
	Natives.Register("simplejvm/http/HttpServer", "isRunning", "()Z", nativeSimpleHttpIsRunning)
}

func nativeHttpServerCreate(frame *Frame) error {
	httpServer = NewHTTPServer()
	return nil
}

func nativeHttpServerStart(frame *Frame) error {
	port := frame.OperandStack.PopInt()
	if httpServer == nil {
		httpServer = NewHTTPServer()
	}
	err := httpServer.Start(int(port))
	if err != nil {
		frame.OperandStack.PushInt(0) // false
	} else {
		frame.OperandStack.PushInt(1) // true
	}
	return nil
}

func nativeHttpServerStop(frame *Frame) error {
	if httpServer != nil {
		httpServer.Stop()
	}
	return nil
}

func nativeHttpServerHandle(frame *Frame) error {
	stack := frame.OperandStack
	statusCode := stack.PopInt()
	response := popString(stack)
	path := popString(stack)
	method := popString(stack)

	if httpServer == nil {
		httpServer = NewHTTPServer()
	}
	httpServer.RegisterHandler(method, path, response, int(statusCode), "text/plain")
	return nil
}

func nativeHttpServerHandleJson(frame *Frame) error {
	stack := frame.OperandStack
	statusCode := stack.PopInt()
	response := popString(stack)
	path := popString(stack)
	method := popString(stack)

	if httpServer == nil {
		httpServer = NewHTTPServer()
	}
	httpServer.RegisterHandler(method, path, response, int(statusCode), "application/json")
	return nil
}

// popString pops a string from the operand stack
func popString(stack *OperandStack) string {
	ref := stack.PopRef()
	if s, ok := ref.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", ref)
}

// =============== SimpleHttp - Integer-based API for current JVM limitations ===============

// Simple mapping: method 1=GET, 2=POST, 3=PUT, 4=DELETE
// Simple mapping: path 1=/, 2=/api/hello, 3=/api/data, 4=/api/users, 5=/health

var simplePathMap = map[int32]string{
	1: "/",
	2: "/api/hello",
	3: "/api/data",
	4: "/api/users",
	5: "/health",
}

var simpleMethodMap = map[int32]string{
	1: "GET",
	2: "POST",
	3: "PUT",
	4: "DELETE",
	0: "*", // wildcard
}

func nativeSimpleHttpStart(frame *Frame) error {
	port := frame.OperandStack.PopInt()
	if httpServer == nil {
		httpServer = NewHTTPServer()
	}
	err := httpServer.Start(int(port))
	if err != nil {
		frame.OperandStack.PushInt(0)
	} else {
		frame.OperandStack.PushInt(1)
	}
	return nil
}

func nativeSimpleHttpStop(frame *Frame) error {
	if httpServer != nil {
		httpServer.Stop()
		httpServer = nil
	}
	return nil
}

func nativeSimpleHttpAddRoute(frame *Frame) error {
	stack := frame.OperandStack
	statusCode := stack.PopInt()
	responseRef := stack.PopRef()
	pathId := stack.PopInt()
	methodId := stack.PopInt()

	method := simpleMethodMap[methodId]
	path := simplePathMap[pathId]

	response := ""
	if s, ok := responseRef.(string); ok {
		response = s
	} else {
		response = fmt.Sprintf("%v", responseRef)
	}

	if httpServer == nil {
		httpServer = NewHTTPServer()
	}

	// Determine content type based on response
	contentType := "text/plain"
	if strings.HasPrefix(response, "{") || strings.HasPrefix(response, "[") {
		contentType = "application/json"
	}

	httpServer.RegisterHandler(method, path, response, int(statusCode), contentType)
	return nil
}

func nativeSimpleHttpIsRunning(frame *Frame) error {
	if httpServer != nil {
		httpServer.mu.RLock()
		running := httpServer.running
		httpServer.mu.RUnlock()
		if running {
			frame.OperandStack.PushInt(1)
			return nil
		}
	}
	frame.OperandStack.PushInt(0)
	return nil
}
