package main

import (
    "encoding/json"
    "log"
    "net/http"
    "airbnb-analytics/internal/service"
    "github.com/gorilla/mux"
)

func main() {
    roomService := service.NewRoomService()
    router := mux.NewRouter()

    // CORS middleware
    router.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    })

    // Room analytics endpoint
    router.HandleFunc("/{roomId}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        roomID := vars["roomId"]

        analytics, err := roomService.GetRoomAnalytics(roomID)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(analytics)
    }).Methods("GET", "OPTIONS")

    log.Printf("Server starting on port 8080...")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatal(err)
    }
}
