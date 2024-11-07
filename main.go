package main

import (
	"log"
	"net/http"

	"github.com/amankumar94728/fealtyx-student-api/internal/api"
	"github.com/amankumar94728/fealtyx-student-api/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	store := storage.NewStorage()
	apiHandler := api.NewAPI(store)

	r := mux.NewRouter()

	r.HandleFunc("/students", apiHandler.CreateStudent).Methods("POST")
	r.HandleFunc("/students", apiHandler.GetAllStudents).Methods("GET")
	r.HandleFunc("/students/{id}", apiHandler.GetStudentByID).Methods("GET")
	r.HandleFunc("/students/{id}", apiHandler.UpdateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", apiHandler.DeleteStudent).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", apiHandler.GenerateStudentSummary).Methods("GET")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
