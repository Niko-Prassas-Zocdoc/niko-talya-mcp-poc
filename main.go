package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// QuestionRequest represents the incoming question request
type QuestionRequest struct {
	Question string `json:"question"`
}

// AnswerResponse represents the response to a question
type AnswerResponse struct {
	Answer string `json:"answer"`
}

func main() {
	// Define HTTP endpoints
	http.HandleFunc("/ask", handleQuestion)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

// handleQuestion processes incoming questions
func handleQuestion(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req QuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// For now, just return a placeholder response
	response := AnswerResponse{
		Answer: "This is a placeholder answer. Your question was: " + req.Question,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
