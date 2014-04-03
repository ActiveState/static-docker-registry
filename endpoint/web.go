package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)


func main() {
	port := os.Getenv("PORT")
	addr := ":" + port
	cdnUrl := os.Getenv("CDN_URL")

	if cdnUrl == "" {
		fmt.Printf("ERROR: CDN_URL must be specified")
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		path := r.URL.Path

		// All request URLs must begin with /v1/ 
		if !strings.HasPrefix(path, "/v1/") {
			w.WriteHeader(404)
			return
		}
		path = path[4:]

		// Serve /v1/_ping
		if path == "_ping" {
			fmt.Fprint(w, "Endpoint app")
			return
		}

		url := cdnUrl + path
		fmt.Printf("[REDIRECT] %v\n", url)
		http.Redirect(w, r, url, http.StatusFound)
	})

	fmt.Printf("Listening at %v\n", addr)
	http.ListenAndServe(addr, nil)
}
