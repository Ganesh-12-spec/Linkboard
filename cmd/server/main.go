package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello fro linkBoard")
	})

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
