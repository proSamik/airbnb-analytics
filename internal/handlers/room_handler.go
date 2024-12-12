package handlers

import (
	"airbnb-analytics/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// HandleRoomAnalytics creates a handler for room analytics requests.
// Parameters:
//   - roomService *service.RoomService: Service for processing room analytics
//
// Returns:
//   - http.HandlerFunc: Handler function for room analytics endpoint
func HandleRoomAnalytics(roomService *service.RoomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID := vars["roomId"]

		// Validate room ID
		if strings.TrimSpace(roomID) == "" {
			handleError(w, "room ID is required", http.StatusBadRequest)
			return
		}

		analytics, err := roomService.GetRoomAnalytics(roomID)
		if err != nil {
			if err.Error() == "room not found" {
				handleError(w, "room not found", http.StatusNotFound)
				return
			}
			handleError(w, "failed to fetch room analytics", http.StatusInternalServerError)
			return
		}

		sendJSONResponse(w, analytics)
	}
}

// HandleGetAllRooms creates a handler for retrieving all available room IDs.
// Parameters:
//   - roomService *service.RoomService: Service for room operations
//
// Returns:
//   - http.HandlerFunc: Handler function for getting all rooms endpoint
func HandleGetAllRooms(roomService *service.RoomService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := roomService.GetAllRooms()
		if err != nil {
			handleError(w, "failed to fetch rooms", http.StatusInternalServerError)
			return
		}

		if len(rooms) == 0 {
			sendJSONResponse(w, map[string][]string{"rooms": {}})
			return
		}

		sendJSONResponse(w, map[string][]string{"rooms": rooms})
	}
}

// handleError processes and sends an error response to the client.
// Parameters:
//   - w http.ResponseWriter: Response writer to send error
//   - message string: Error message to send
//   - statusCode int: HTTP status code
func handleError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": message}); encodeErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// sendJSONResponse sends a JSON response to the client.
// Parameters:
//   - w http.ResponseWriter: Response writer to send JSON
//   - data interface{}: Data to encode as JSON
func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
