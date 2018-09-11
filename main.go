package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/imt2681-assig1/igcinfo"
)

const rootPath = "igcinfo"

func handleRoutes(w http.ResponseWriter, r *http.Request) {
	// Get paths
	routes := strings.Split(r.URL.Path[1:], "/")[2:]

	// GET /api, send info about the API
	if routes[0] == "" {
		igcinfo.APIInfo(w, r)
		return
	}

	if routes[0] == "igc" {
		igcRoute := routes[1:]

		if r.Method == http.MethodGet {
			// GET /api/igc
			if len(igcRoute) == 0 || igcRoute[0] == "" {
				igcinfo.GetAllTracks(w, r)
				return
			}
		}

		if r.Method == http.MethodPost {
			// POST /api/igc
			if len(igcRoute) == 0 || igcRoute[0] == "" {
				igcinfo.RegisterTrack(w, r)
				return
			}
		}
	}

	http.NotFound(w, r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT not set")
	}

	http.HandleFunc("/igcinfo/api/", handleRoutes)

	// Start listen
	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
