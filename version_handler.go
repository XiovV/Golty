package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type version struct {
	VersionNumber string
}

func HandleGetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	logrus.Info("returning version number")
	json.NewEncoder(w).Encode(version{VersionNumber: VERSION})
}