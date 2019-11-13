package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func initLogFile() {
	logFile, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("Error reading log file: ", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)

	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(mw)
}

func initChannelsDatabase() {
	_, err := os.Stat("./config/channels.json")
	if os.IsNotExist(err) {
		f, _ := os.Create("./config/channels.json")
		defer f.Close()
		s, _ := f.WriteString("[]")
		log.Info("initiated channels.json: ", s)
		f.Sync()
	}
}

func init() {
	initLogFile()
	initChannelsDatabase()
}

func uploadChecker() {
	for {
		time.Sleep(5 * time.Hour)

		// go CheckAll()
		log.Info("Upload Checker running...")
	}
}

func main() {
	log.Info("Server running on port 8080")

	go uploadChecker()

	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/", HandleIndex).Methods("GET")
	r.HandleFunc("/logs", HandleLogs).Methods("GET")
	r.HandleFunc("/api/get-channels", HandleGetChannels).Methods("GET")
	r.HandleFunc("/api/add-channel", HandleAddChannel).Methods("POST")
	r.HandleFunc("/api/check-channel", HandleCheckChannel).Methods("POST")
	r.HandleFunc("/api/check-all", HandleCheckAll).Methods("GET")
	r.HandleFunc("/api/delete-channel", HandleDeleteChannel).Methods("POST")

	r.HandleFunc("/static/app.js", ServeJS).Methods("GET")

	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r))
}
