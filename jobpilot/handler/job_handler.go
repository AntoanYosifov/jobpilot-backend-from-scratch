package handler

import (
	"encoding/json"
	"jobpilot/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	jobs   = []model.Job{}
	nextID = 1
	mu     sync.Mutex
)

func JobsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	switch r.Method {
	case http.MethodGet:
		GetJobs(w, r)
	case http.MethodPost:
		CreateJob(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func JobByIdHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/jobs/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	updateJobFullReplace(w, r, id)
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(jobs)
	if err != nil {
		return
	}
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
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

func updateJobFullReplace(w http.ResponseWriter, r *http.Request, id int) {
	var updated model.Job
	err := json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	if updated.Title == "" || updated.Company == "" || updated.Status == "" || updated.Date == "" {
		http.Error(w, "All fields are required for PUT", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range jobs {
		if jobs[i].ID == id {
			updated.ID = id
			jobs[i] = updated

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(jobs[i])
			return
		}
	}

	http.Error(w, "Job not found", http.StatusNotFound)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
