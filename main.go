package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	port   string
	dlRoot string
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

func initDatabase() {
	log.Info("Initiating db")
	databases := [3]string{"./config/channels.json", "./config/playlists.json", "./config/videos.json"}
	for _, database := range databases {
		createDir(database)
	}
}

func init() {
	flag.StringVar(&port, "port", "8080", "port to listen")
	flag.StringVar(&dlRoot, "dl", "./downloads", "root download folder")
}

func uploadCheckerChannels() {
	interval, err := GetCheckingInterval(TypeChannel)
	if err != nil {
		log.Errorf("uploadCheckerChannels: %v", err)
	}
	if interval == 0 {
		time.Sleep(5 * time.Minute)
		uploadCheckerChannels()
	} else if interval != 0 {
		for {
			if interval != 0 {
				time.Sleep(time.Duration(interval) * time.Minute)
				go CheckAll("channels")
				log.Infof("upload Checker for channels running every %v minutes", interval)
			}
		}
	}
}

func uploadCheckerPlaylists() {
	interval, err := GetCheckingInterval(TypePlaylist)
	if err != nil {
		log.Errorf("uploadCheckerPlaylists: %v", err)
	}
	if interval == 0 {
		time.Sleep(5 * time.Minute)
		uploadCheckerPlaylists()
	} else if interval != 0 {
		for {
			if interval != 0 {
				<-time.After(time.Duration(interval) * time.Minute)
				go CheckAll("playlists")
				log.Infof("upload Checker for playlists running every %v minutes", interval)
			}
		}
	}
}

func main() {
	flag.Parse()
	log.Info("server running on port " + port)
	go uploadCheckerChannels()
	go uploadCheckerPlaylists()

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

	r.HandleFunc("/api/get-checking-interval", HandleGetCheckingInterval).Methods("POST")
	r.HandleFunc("/api/update-checking-interval", HandleUpdateCheckingInterval).Methods("POST")

	r.HandleFunc("/api/get-videos", HandleGetVideos).Methods("GET")
	r.HandleFunc("/api/download-video", HandleDownloadVideo).Methods("POST")

	r.HandleFunc("/api/version", HandleGetVersion).Methods("GET")

	r.HandleFunc("/favicon.ico", ServeFavicon).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(r)))
}
