package handlers

import (
	"airbnb-analytics/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
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

		analytics, err := roomService.GetRoomAnalytics(roomID)
		if err != nil {
			handleError(w, err)
			return
		}

		sendJSONResponse(w, analytics)
	}
}

// handleError processes and sends an error response to the client.
// Parameters:
//   - w http.ResponseWriter: Response writer to send error
//   - err error: Error to process and send
func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); encodeErr != nil {
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
