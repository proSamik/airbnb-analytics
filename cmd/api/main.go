package main

import (
	"airbnb-analytics/internal/handlers"
	"airbnb-analytics/internal/middleware"
	"airbnb-analytics/internal/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// main initializes and starts the HTTP server with configured routes
// and middleware on port 8080.
func main() {
	// Initialize services
	roomService := service.NewRoomService()

	// Initialize router
	router := setupRouter(roomService)

	// Start server
	startServer(router)
}

// setupRouter initializes and configures the router with all middleware and handlers.
// Parameters:
//   - roomService *service.RoomService: Service handling room analytics logic
//
// Returns:
//   - *mux.Router: Configured router instance
func setupRouter(roomService *service.RoomService) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.CORS)

	// Register routes
	registerRoutes(router, roomService)

	return router
}

// registerRoutes sets up all API endpoints for the application.
// Parameters:
//   - router *mux.Router: Router to register routes on
//   - roomService *service.RoomService: Service for room analytics
func registerRoutes(router *mux.Router, roomService *service.RoomService) {
	// Room analytics endpoint
	router.HandleFunc("/{roomId}",
		handlers.HandleRoomAnalytics(roomService),
	).Methods("GET", "OPTIONS")
}

// startServer starts the HTTP server on the specified port.
// If the server fails to start, it logs the error and exits.
// Parameters:
//   - router *mux.Router: Configured router to use for the server
func startServer(router *mux.Router) {
	const port = ":8080"
	log.Printf("Server starting on port %s...", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
