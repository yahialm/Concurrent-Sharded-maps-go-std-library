package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

func TestConcurrentPostRequests(t *testing.T) {
	const numRequests = 200
	const baseURL = "http://localhost:8080/api/" // Replace with your API base URL

	var wg sync.WaitGroup

	for i := 1; i <= numRequests; i++ {
		wg.Add(1)
		go func(key string, value any) {
			defer wg.Done()

			// Prepare JSON payload
			payload, err := json.Marshal(map[string]interface{}{"key": key, "value": value})
			if err != nil {
				t.Errorf("Error marshalling JSON: %v", err)
				return
			}

			// Send POST request
			resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				t.Errorf("POST request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			// Check response status code
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code: %d", resp.StatusCode)
			}
		}(fmt.Sprintf("key%d", i), i) // Use loop variables to create string keys and integer values
	}

	// Wait for all requests to finish
	wg.Wait()
}
