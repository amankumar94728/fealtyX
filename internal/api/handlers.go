package api

import (
    "bytes"
    "io"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/amankumar94728/fealtyx-student-api/internal/models"
    "github.com/amankumar94728/fealtyx-student-api/internal/storage"
)

type API struct {
    storage *storage.Storage
}

type GenerateRequest struct { 
    Model string `json:"model"` 
    Prompt string `json:"prompt"` 
    Stream bool `json:"stream"` 
    } 
type GenerateResponse struct {
     Model string `json:"model"` 
     CreatedAt string `json:"created_at"` 
     Response string `json:"response"` 
     Done bool `json:"done"` 
     Context []int `json:"context"` 
     TotalDuration int64 `json:"total_duration"` 
     LoadDuration int64 `json:"load_duration"` 
     PromptEvalCount int `json:"prompt_eval_count"` 
     PromptEvalDuration int64 `json:"prompt_eval_duration"` 
     EvalCount int `json:"eval_count"` 
     EvalDuration int64 `json:"eval_duration"` }

func NewAPI(storage *storage.Storage) *API {
    return &API{storage: storage}
}

func (a *API) CreateStudent(w http.ResponseWriter, r *http.Request) {
    var student models.Student
    if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := validateStudent(student); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    createdStudent, err := a.storage.Create(student)
    if err != nil {
        http.Error(w, "Failed to create student", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdStudent)
}

func (a *API) GetAllStudents(w http.ResponseWriter, r *http.Request) {
    students := a.storage.GetAll()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(students)
}

func (a *API) GetStudentByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    student, err := a.storage.GetByID(id)
    if err != nil {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(student)
}

func (a *API) UpdateStudent(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    var student models.Student
    if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if err := validateStudent(student); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    updatedStudent, err := a.storage.Update(id, student)
    if err != nil {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedStudent)
}

func (a *API) DeleteStudent(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    if err := a.storage.Delete(id); err != nil {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (a *API) GenerateStudentSummary(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    student, err := a.storage.GetByID(id)
    if err != nil {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

    summary, err := generateSummaryWithOllama(student)
    if err != nil {
        http.Error(w, "Failed to generate summary", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}

func validateStudent(student models.Student) error {
    if student.Name == "" {
        return fmt.Errorf("name is required")
    }
    if student.Age <= 0 {
        return fmt.Errorf("age must be positive")
    }
    if student.Email == "" {
        return fmt.Errorf("email is required")
    }
    // You can add more validation rules here
    return nil
}

func generateSummaryWithOllama(student models.Student) (string, error) {
    prompt := fmt.Sprintf("Generate a brief summary for a student named %s, who is %d years old and has the email %s.", student.Name, student.Age, student.Email)
    url := "http://localhost:12345/api/generate"
    requestBody := GenerateRequest{ 
        Model: "llama3.2:1b", 
        Prompt: prompt, 
        Stream: false, 
    } 
    jsonStr, err := json.Marshal(requestBody)
    if err != nil { 
        return "", fmt.Errorf("error marshalling request body: %v", err) 
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    if err != nil { 
        return "", fmt.Errorf("error creating request: %v", err) 
    } 
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{} 
    resp, err := client.Do(req) 
    if err != nil { 
        return "", fmt.Errorf("error making request: %v", err) 
    } 
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response: %v", err)
    }

    var generateResponse GenerateResponse 
    err = json.Unmarshal(body, &generateResponse) 
    if err != nil { 
        return "", fmt.Errorf("error unmarshalling response: %v", err) 
    } 
    return generateResponse.Response, nil

}