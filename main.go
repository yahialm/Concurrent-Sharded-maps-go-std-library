package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {

	// To handle concurrent requests, each request to the api should run on its own goroutine
	// That's what I am doing here

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {

		var wg sync.WaitGroup = sync.WaitGroup{} 

		wg.Add(1)
		switch r.Method {
		case http.MethodGet:
			go handleGet(w, r, &wg)
		case http.MethodPost:
			go handlePost(w, r, &wg)
		case http.MethodDelete:
			go handleDelete(w, r, &wg)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method %s not allowed\n", r.Method)
		}
		wg.Wait() // wait for go routines to finish
	})

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}