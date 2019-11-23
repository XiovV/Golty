package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
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
	interval, err := GetCheckingInterval()
	if err != nil {
		log.Errorf("uploadChecker: %s", err)
	}
	if interval == 0 {
		time.Sleep(5 * time.Minute)
		uploadChecker()
	} else if interval != 0 {
		for {
			if interval != 0 {
				time.Sleep(time.Duration(interval) * time.Minute)
				go CheckAll("channels")
				log.Infof("upload Checker running every %v minutes", interval)
			}
		}
	}
}

func main() {
	log.Info("server running on port 8080")
	
	go uploadChecker()

	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	r.HandleFunc("/", HandleIndex).Methods("GET")
	r.HandleFunc("/logs", HandleLogs).Methods("GET")
	r.HandleFunc("/playlists", HandlePlaylists).Methods("GET")
	r.HandleFunc("/videos", HandleVideos).Methods("GET")

	r.HandleFunc("/api/add", HandleAddTarget).Methods("POST")
	r.HandleFunc("/api/check", HandleCheckTarget).Methods("POST")
	r.HandleFunc("/api/delete", HandleDeleteTarget).Methods("POST")
	r.HandleFunc("/api/get", HandleGetTargets).Methods("POST")
	r.HandleFunc("/api/check-all", HandleCheckAllTargets).Methods("POST")

	r.HandleFunc("/api/update-checking-interval", HandleUpdateCheckingInterval).Methods("POST")


	r.HandleFunc("/api/get-videos", HandleGetVideos).Methods("GET")
	r.HandleFunc("/api/download-video", HandleDownloadVideo).Methods("POST")

	r.HandleFunc("/api/version", HandleGetVersion).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r))
}
