package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

func TestConcurrentDeleteRequests(t *testing.T) {
	const numRequests = 200
	const baseURL = "http://localhost:8080/api/"

	var wg sync.WaitGroup

	client := &http.Client{} // Create an HTTP client

	for i := 1; i <= numRequests; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()

			// Create a DELETE request
			req, err := http.NewRequest("DELETE", baseURL + key, nil)
			if err != nil {
				t.Errorf("DELETE request failed: %v", err)
				return
			}

			// Send DELETE request
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("DELETE request failed: %v", err)
				return
			}

			defer resp.Body.Close()

			// Check response status code
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code: %d", resp.StatusCode)
			}
		}(fmt.Sprintf("key%d", i))
	}

	// Wait for all requests to finish
	wg.Wait()

}