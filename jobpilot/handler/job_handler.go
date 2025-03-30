package handler

import (
	"encoding/json"
	"jobpilot/model"
	"log"
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

func CreateJob(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.Header().Set("Content-Type", "application/json")

	var newJobs []model.Job
	err := json.NewDecoder(r.Body).Decode(&newJobs)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	mu.Lock()
	for i := range newJobs {
		if newJobs[i].Title == "" || newJobs[i].Company == "" {
			mu.Unlock()
			http.Error(w, "Title and Company are required for all jobs", http.StatusBadRequest)
			return
		}
		newJobs[i].ID = nextID
		nextID++
		jobs = append(jobs, newJobs[i])
	}
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newJobs)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
