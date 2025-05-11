package main

//5432 Sever
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/cors"
)

// Define request and response structures
type RequestData struct {
	Message string `json:"message"`
}

type ResponseData struct {
	Reply string `json:"reply"`
}

// Handler function to process user messages
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var request RequestData
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create response based on user input
	var response ResponseData
	switch request.Message {
	case "hello":
		response.Reply = "Hello! How can I help you?"
	case "bye":
		response.Reply = "Goodbye! Have a nice day!"
	default:
		response.Reply = "I don't understand that message."
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main0() {
	// Create a new ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", handler)

	// Configure CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow React frontend
		AllowedMethods:   []string{"POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", corsHandler)
}
