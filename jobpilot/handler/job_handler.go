package handler

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"jobpilot/db"
	"jobpilot/model"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func JobsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	switch r.Method {
	case http.MethodGet:
		getJobs(w, r)
	case http.MethodPost:
		createJobs(w, r)
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

	switch r.Method {
	case http.MethodPut:
		updateJobFullReplace(w, r, id)
	case http.MethodPatch:
		updatedJobPartial(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jobs []model.Job
	result := db.DB.Find(&jobs)

	if result.Error != nil {
		log.Printf("[Error] GetJobs DB failure: %v", result.Error)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(jobs)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func createJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newJobs []model.Job
	err := json.NewDecoder(r.Body).Decode(&newJobs)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	for _, job := range newJobs {
		if job.Title == "" || job.Company == "" {
			http.Error(w, "Title and Company are required for all jobs", http.StatusBadRequest)
			return
		}
	}

	result := db.DB.Create(&newJobs)
	if result.Error != nil {
		log.Printf("[ERROR] Failed to insert jobs: %v", result.Error)
		http.Error(w, "Failed to save jobs", http.StatusInternalServerError)
		return
	}

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

	var existing model.Job
	result := db.DB.First(&existing, id)
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}

		log.Printf("[ERROR] DB failed to find job ID %d: %v", id, result.Error)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	existing.Title = updated.Title
	existing.Company = updated.Company
	existing.Status = updated.Status
	existing.Date = updated.Date

	saveResult := db.DB.Save(&existing)
	if saveResult.Error != nil {
		log.Printf("[ERROR] Failed to update job ID %d: %v", id, saveResult.Error)
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(existing)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
	}

}

func updatedJobPartial(w http.ResponseWriter, r *http.Request, id int) {
	var updates model.Job
	err := json.NewDecoder(r.Body).Decode(&updates)

	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	var existing model.Job
	result := db.DB.First(&existing, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Job not found", http.StatusNotFound)
			return
		}
		log.Printf("[ERROR] Failed to fetch job ID %d: %v", id, result.Error)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Company != "" {
		existing.Company = updates.Company
	}
	if updates.Status != "" {
		existing.Status = updates.Status
	}
	if updates.Date != "" {
		existing.Date = updates.Date
	}

	saveResult := db.DB.Save(&existing)
	if saveResult.Error != nil {
		log.Printf("[ERROR] Failed to update job ID %d: %v", id, saveResult.Error)
		http.Error(w, "Failed to update job", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(existing)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
