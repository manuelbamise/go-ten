package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Constants for default values
const (
	defaultPort        = "8080"
	serverReadTimeout  = 15 * time.Second
	serverWriteTimeout = 15 * time.Second
	serverIdleTimeout  = 60 * time.Second
	shutdownTimeout    = 30 * time.Second
)

// HTTP status codes
const (
	statusOK            = http.StatusOK
	statusNotFound      = http.StatusNotFound
	statusInternalError = http.StatusInternalServerError
)

func main() {
	// Load configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/ping", pingHandler)

	// Create not found handler
	notFoundMux := http.NotFoundHandler()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the path matches any registered route
		if r.URL.Path != "/health" && r.URL.Path != "/api/v1/ping" {
			notFoundHandler(w, r)
			return
		}
		// If it's a registered route but method doesn't match, let the default handler handle it
		notFoundMux.ServeHTTP(w, r)
	})

	// Apply middleware chain
	handler := applyMiddleware(mux)

	// Setup server with timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}

	// Start server with graceful shutdown
	startServer(server, port)
}

// applyMiddleware applies middleware chain to the handler
func applyMiddleware(handler http.Handler) http.Handler {
	handler = loggingMiddleware(handler)
	handler = corsMiddleware(handler)
	handler = recoveryMiddleware(handler)
	return handler
}

// loggingMiddleware logs method, path, status code, and duration
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: statusOK}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log request details
		duration := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// corsMiddleware adds basic CORS headers for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(statusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// recoveryMiddleware catches panics and returns 500 with error message
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				writeJSON(w, statusInternalError, map[string]string{"error": "Internal server error"})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// writeJSON writes JSON responses with proper Content-Type header
func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// healthHandler returns health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	writeJSON(w, statusOK, response)
}

// pingHandler returns simple ping endpoint
func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	response := map[string]string{"message": "pong"}
	writeJSON(w, statusOK, response)
}

// notFoundHandler handles undefined routes
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"error": "route not found"}
	writeJSON(w, statusNotFound, response)
}

// startServer starts the server with graceful shutdown
func startServer(server *http.Server, port string) {
	// Channel to listen for errors
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either error or shutdown signal
	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)
	case sig := <-shutdown:
		log.Printf("Received signal %v, shutting down server...", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		} else {
			log.Printf("Server stopped gracefully")
		}
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
