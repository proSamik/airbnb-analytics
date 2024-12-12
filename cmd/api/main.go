package main

import (
	"airbnb-analytics/internal/database"
	"airbnb-analytics/internal/handlers"
	"airbnb-analytics/internal/middleware"
	"airbnb-analytics/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// main initializes and starts the HTTP server.
// It performs the following operations in order:
// 1. Loads environment variables from .env file (if exists)
// 2. Initializes database connection
// 3. Sets up services and routing
// 4. Starts HTTP server on configured port
//
// The server will exit if any initialization step fails.
func main() {
	// Load environment variables, don't fail if .env doesn't exist
	if err := godotenv.Load(); err != nil {
		// Just log warning since .env is optional in production
		log.Fatalf("Warning: .env file not found: %v", err)
	}

	// Initialize database connection
	if err := database.InitDB(); err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Initialize services
	roomService := service.NewRoomService()

	// Initialize router
	router := setupRouter(roomService)

	// Start server
	startServer(router)
}

// setupRouter initializes and configures the HTTP router.
// It sets up middleware and routes for the application.
//
// Parameters:
//   - roomService *service.RoomService: Service instance handling room analytics logic
//
// Returns:
//   - *mux.Router: Configured router instance ready for use
//
// The router is configured with CORS middleware and all application routes.
func setupRouter(roomService *service.RoomService) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.CORS)

	// Register routes
	registerRoutes(router, roomService)

	return router
}

// registerRoutes configures all API endpoints for the application.
// It sets up the following routes:
// - GET /health: Health check endpoint
// - GET /rooms: Returns list of all available room IDs
// - GET /{roomId}: Returns analytics for a specific room
//
// Parameters:
//   - router *mux.Router: Router instance to register routes on
//   - roomService *service.RoomService: Service handling room analytics operations
//
// Each route supports both GET and OPTIONS methods for CORS compatibility.
func registerRoutes(router *mux.Router, roomService *service.RoomService) {
	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"message": "Service is running",
		}); err != nil {
			log.Printf("Error encoding health check response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Get all available room IDs
	router.HandleFunc("/rooms",
		handlers.HandleGetAllRooms(roomService),
	).Methods("GET", "OPTIONS")

	// Get analytics for a specific room
	router.HandleFunc("/{roomId}",
		handlers.HandleRoomAnalytics(roomService),
	).Methods("GET", "OPTIONS")
}

// startServer initializes and starts the HTTP server.
// It listens on the configured port and handles incoming HTTP requests.
// The port is determined from environment variable PORT, defaults to 8080.
//
// Parameters:
//   - router *mux.Router: Configured router to handle incoming requests
//
// The function will log fatal error and exit if server fails to start.
func startServer(router *mux.Router) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not set
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
