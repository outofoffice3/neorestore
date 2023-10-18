package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Create a struct to represent the JSON data
type RequestData struct {
	Prefix string `json:"prefix"`
}

func main() {
	// Define a handler function to respond to requests.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a variable to hold the JSON data
		var requestData RequestData

		// Parse the JSON request body
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&requestData)
		if err != nil {
			http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Access the data in the variable
		prefix := requestData.Prefix
		fmt.Println("Received Prefix:", prefix)

		// Respond with a message that includes the received prefix
		response := fmt.Sprintf("Received Prefix: %s", prefix)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, response)
	})

	port := 8080 // Specify the port you want to listen on.

	// Start the HTTP server and listen on the specified port.
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
