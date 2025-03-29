package handler

import (
	"encoding/json"
	"jobpilot/model"
	"net/http"
	"sync"
)

var (
	jobs   = []model.Job{}
	nextID = 1
	mu     sync.Mutex
)

func GetJobs(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(jobs)
	if err != nil {
		return
	}
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
