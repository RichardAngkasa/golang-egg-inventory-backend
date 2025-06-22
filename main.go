package main

import (
	"fmt"
	"log"
	"net/http"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	db := NewDatabase()
	auth := NewAuth("my-super-secret-key")
	handlers := NewHandlers(db, auth)

	http.HandleFunc("/login", enableCORS(handlers.Login))
	http.HandleFunc("/egg-racks", enableCORS(handlers.GetAllEggRacks))
	http.HandleFunc("/egg-racks/bulk", enableCORS(handlers.CreateBulkEggRacks))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
