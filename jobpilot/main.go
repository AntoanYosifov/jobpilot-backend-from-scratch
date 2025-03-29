package main

import (
	"fmt"
	"jobpilot/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetJobs(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println("🚀 JobPilot backend running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
