package main

import (
	"fmt"
	"jobpilot/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/jobs", handler.JobsHandler)
	http.HandleFunc("/jobs/", handler.JobByIdHandler)
	fmt.Println("🚀 JobPilot backend running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
