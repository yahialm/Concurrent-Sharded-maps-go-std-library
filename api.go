package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var sMap shardedMap = *NewShardedMap()

// Getting the key from URLs
func getKeyFromURL(urlPath string) string {
    parts := strings.Split(urlPath, "/")
    return parts[len(parts)-1]
}

// method "GET" handler
func handleGet(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {

	defer wg.Done() // decrement the wg counter

    key := getKeyFromURL(r.URL.Path)
    value, err := sMap.get(key)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "Key '%s' not found\n", key)
        return
    }
    
    // Create a map to hold key-value pair
    kv := map[string]interface{}{
        key: value,
    }

    // Convert map to JSON
    jsonData, err := json.Marshal(kv)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error marshalling data: %v\n", err)
        return
    }

    // Set the content type and write the JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
    fmt.Printf("GET request handled successfully for key '%s'\n", key)
}

func handlePost(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {

	defer wg.Done() // // decrement the wg counter

    var data struct {
        Key   string `json:"key"`
        Value any `json:"value"`
    }
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, "Error decoding request body: %v\n", err)
        return
    }
    _, err := sMap.store(data.Key, data.Value)

	sErr := fmt.Sprintln(err)

	if err != nil {
		fmt.Fprintf(w, "Error: %s",sErr)
		return
	}
    fmt.Fprintf(w, "Stored key '%s' with value '%s'\n", data.Key, data.Value)

    // Log a simple message to the terminal
    fmt.Printf("POST request handled successfully for key '%s'\n", data.Key)
}

func handleDelete(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {

	defer wg.Done() // decrement the wg counter

    key := getKeyFromURL(r.URL.Path)
    err := sMap.delete(key)

	if err != nil {
		fmt.Println(err) 
		return 
	}

    fmt.Fprintf(w, "Deleted key '%s'\n", key)
}