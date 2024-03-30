package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {

	targets := map[string]string{
		"python": "http://host.docker.internal:7080",
		"node":   "http://host.docker.internal:8080",
	}

	// Open log file for writing
	logFile, err := os.OpenFile("proxy.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Create logger that writes to log file
	logger := log.New(logFile, "PROXY ", log.LstdFlags)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		splitPaths := strings.SplitN(r.URL.Path, "/", 3)
		if len(splitPaths) < 2 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			logger.Println("Invalid request")
			return
		}
		namespace := splitPaths[1]
		target, ok := targets[namespace]
		if !ok {
			http.Error(w, "Unknown namespace", http.StatusNotFound)
			logger.Printf("Unknown namespace: %s\n", namespace)
			return
		}

		remote, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			logger.Printf("Error parsing URL: %v\n", err)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		director := proxy.Director
		proxy.Director = func(req *http.Request) {
			director(req)
			logger.Println(req.URL.Path)

			splitPath := strings.SplitN(req.URL.Path, "/", 3)
			logger.Println(splitPath)
			if len(splitPath) > 2 {
				req.URL.Path = "/" + splitPath[2]
			} else {
				req.URL.Path = "/"
			}

			req.Header.Set("Host", req.Host)
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)
			req.Header.Set("X-Forwarded-Proto", req.URL.Scheme)
			req.Header.Set("X-Real-IP", req.RemoteAddr)
		}

		r.URL.Host = remote.Host
		r.URL.Scheme = remote.Scheme
		proxy.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		logger.Fatalf("Server error: %v\n", err)
	}
}
