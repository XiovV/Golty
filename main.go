package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Server running")

	go UploadChecker()

	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/", HandleIndex).Methods("GET")
	r.HandleFunc("/api/get-channels", HandleGetChannels).Methods("GET")
	// r.HandleFunc("/addchannel", HandleAddChannel).Methods("POST")
	r.HandleFunc("/api/check-channel", HandleCheckChannel).Methods("POST")
	r.HandleFunc("/api/check-all", HandleCheckAll).Methods("GET")

	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r))
}
