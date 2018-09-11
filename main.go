package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const rootPath = "igcinfo"

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Get request paths
	paths := strings.Split(r.URL.Path[1:], "/")

	fmt.Printf("Paths len: %d\n", len(paths))
	fmt.Println(paths)

	// If root path is not igcinfo, next path not api, reply with 404
	if len(paths) != 2 || paths[0] != rootPath || paths[1] != "api" {
		http.NotFound(w, r)
		return
	}

	handleRoutes(&w, r, paths[2:])
}

func handleRoutes(w *http.ResponseWriter, r *http.Request, routes []string) {
	// Base api route
	if len(routes) == 0 {

	}
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT not set")
	}

	// Start listen
	http.HandleFunc("/", mainHandler)
	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)

	// Exit if error
	if err != nil {
		panic(err)
	}
}
