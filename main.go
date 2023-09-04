package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

func forwardRequest(request string) (string, error) {
	// Create an HTTP client
	client := &http.Client{}

	// Create an HTTP request with the JSON payload
	req, err := http.NewRequest("POST", "http://localhost:9933", bytes.NewBuffer([]byte(request)))
	if err != nil {
		return "", err
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	limiter := rate.NewLimiter(rate.Limit(1), 1) // Allow 1 requests per 1 second

	// Upgrade the HTTP connection to a WebSocket connection
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// Handle error
		return
	}

	// Read incoming messages from the WebSocket connection
	for {
		// Apply rate limiting
		if !limiter.Allow() {
			fmt.Println("Rate limit exceeded")
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			// Handle error
			break
		}

		// Forward the JSON request to another HTTP endpoint
		response, err := forwardRequest(string(message))
		if err != nil {
			// Handle error
			break
		}

		// Send the response back to the WebSocket client
		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			// Handle error
			break
		}
	}
}

func main() {
	fmt.Println("Hello, World!")
	http.HandleFunc("/", websocketHandler)
	http.ListenAndServe(":8080", nil) // Replace 8080 with the desired port number
}
