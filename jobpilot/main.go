package main

import (
	"fmt"
	"jobpilot/db"
	"jobpilot/handler"
	"jobpilot/model"
	"log"
	"net/http"
)

func main() {
	db.Connect()
	err := db.DB.AutoMigrate(&model.Job{})

	if err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	http.HandleFunc("/jobs", handler.JobsHandler)
	http.HandleFunc("/jobs/", handler.JobByIdHandler)
	fmt.Println("🚀 JobPilot backend running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
