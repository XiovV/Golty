package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func uploadChecker() {
	for {
		time.Sleep(20 * time.Second)

		// go CheckAll()
		fmt.Println("Upload Checker running...")
	}
}

func main() {
	fmt.Println("Server running on port 8080")

	go uploadChecker()

	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/", HandleIndex).Methods("GET")
	r.HandleFunc("/api/get-channels", HandleGetChannels).Methods("GET")
	r.HandleFunc("/api/add-channel", HandleAddChannel).Methods("POST")
	r.HandleFunc("/api/check-channel", HandleCheckChannel).Methods("POST")
	r.HandleFunc("/api/check-all", HandleCheckAll).Methods("GET")
	r.HandleFunc("/api/delete-channel", HandleDeleteChannel).Methods("POST")

	r.HandleFunc("/static/app.js", ServeJS).Methods("GET")

	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r))
}
