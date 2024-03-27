package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

func TestConcurrentGetRequests(t *testing.T) {
	const numRequests = 200
	const baseURL = "http://localhost:8080/api/" // Replace with your API base URL

	var wg sync.WaitGroup

	// Send concurrent GET requests
	for i := 1; i <= numRequests; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()

			// Send GET request to the server
			resp, err := http.Get(baseURL + key)
			if err != nil {
				t.Errorf("GET request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			// Check response status code
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code: %d", resp.StatusCode)
				return
			}

			// Process response if needed
		}(fmt.Sprintf("key%d", i))
	}

	// Wait for all requests to finish
	wg.Wait()
}
