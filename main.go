package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	counter int
	mu      sync.Mutex
)

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	counter++
	mu.Unlock()

	response := map[string]int{"counter": counter}
	json.NewEncoder(w).Encode(response)
}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"time": time.Now().UTC().Format(time.RFC3339)}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/v1/counter", CounterHandler)
	http.HandleFunc("/v1/time", TimeHandler)

	fmt.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
