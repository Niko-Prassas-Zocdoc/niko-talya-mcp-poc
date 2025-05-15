package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// MCP protocol structures
type MCPToolCallRequest struct {
	ToolCalls []ToolCall `json:"tool_calls"`
}

type ToolCall struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Parameters json.RawMessage `json:"parameters"`
}

type MCPToolCallResponse struct {
	Responses []ToolResponse `json:"responses"`
}

type ToolResponse struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// QuestionRequest represents the incoming question request
type QuestionRequest struct {
	Question string `json:"question"`
}

// AnswerResponse represents the response to a question
type AnswerResponse struct {
	Answer string `json:"answer"`
}

// Tools declaration for MCP
type MCPTools struct {
	Tools []MCPTool `json:"tools"`
}

type MCPTool struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Parameters  MCPSchema `json:"parameters"`
}

type MCPSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

func main() {
	// Define HTTP endpoints
	http.HandleFunc("/mcp/v1/execute", handleMCPExecute)
	http.HandleFunc("/mcp/v1/tools", handleMCPTools)
	http.HandleFunc("/ask", handleQuestion) // Keep original endpoint for backward compatibility
	http.HandleFunc("/sse", handleSSE)      // Server-Sent Events endpoint for Cursor

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting MCP server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

// handleSSE handles the SSE endpoint for Cursor
func handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Send initial heartbeat
	fmt.Fprintf(w, "event: ready\ndata: {}\n\n")
	w.(http.Flusher).Flush()

	// Keep connection open with periodic heartbeats
	<-r.Context().Done()
}

// handleMCPTools returns the available tools for this MCP server
func handleMCPTools(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define the tools this MCP provides
	tools := MCPTools{
		Tools: []MCPTool{
			{
				Name:        "ask_question",
				Description: "Ask a question and get an answer",
				Parameters: MCPSchema{
					Type: "object",
					Properties: map[string]Property{
						"question": {
							Type:        "string",
							Description: "The question to ask",
						},
					},
					Required: []string{"question"},
				},
			},
		},
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Send the response
	if err := json.NewEncoder(w).Encode(tools); err != nil {
		http.Error(w, "Error encoding tools: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleMCPExecute processes tool call requests from Cursor
func handleMCPExecute(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var req MCPToolCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Process each tool call
	var responses []ToolResponse
	for _, toolCall := range req.ToolCalls {
		switch toolCall.Name {
		case "ask_question":
			// Parse the question parameter
			var params QuestionRequest
			if err := json.Unmarshal(toolCall.Parameters, &params); err != nil {
				responses = append(responses, ToolResponse{
					Type:    "error",
					Content: "Invalid parameters: " + err.Error(),
				})
				continue
			}

			// Generate a response for the question
			answer := "This is a placeholder answer from the MCP server. Your question was: " + params.Question

			responses = append(responses, ToolResponse{
				Type:    "data",
				Content: map[string]string{"answer": answer},
			})
		default:
			responses = append(responses, ToolResponse{
				Type:    "error",
				Content: "Unknown tool: " + toolCall.Name,
			})
		}
	}

	// Create the response
	response := MCPToolCallResponse{
		Responses: responses,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleQuestion processes incoming questions (legacy endpoint)
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
