package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/imt2681-assig1/igcinfo"
)

const rootPath = "igcinfo"

func handleRoutes(w http.ResponseWriter, r *http.Request) {
	// Get paths
	routes := strings.Split(r.URL.Path[1:], "/")[2:]

	// GET /api, send info about the API
	if r.Method == http.MethodGet && routes[0] == "" {
		igcinfo.APIInfo(w, r)
		return
	}

	if routes[0] == "igc" {
		igcRoutes := routes[1:]
		handleIgcRoutes(w, r, igcRoutes)
		return
	}

	http.NotFound(w, r)
}

func handleIgcRoutes(w http.ResponseWriter, r *http.Request, routes []string) {
	if r.Method == http.MethodGet {
		// GET /api/igc
		if len(routes) == 0 || routes[0] == "" {
			igcinfo.GetAllTracks(w, r)
			return
		}

		// GET /api/igc/{id}
		if len(routes) == 1 || routes[1] == "" {
			id, err := strconv.Atoi(routes[0])
			if err == nil {
				igcinfo.GetTrackByID(w, r, id)
				return
			}
		}

		// GET /api/igc/{id}/{field}
		if len(routes) == 2 || routes[2] == "" {
			id, err := strconv.Atoi(routes[0])
			if err == nil {
				igcinfo.GetTrackField(w, r, id, routes[1])
				return
			}
		}
	}

	if r.Method == http.MethodPost {
		// POST /api/igc
		if len(routes) == 0 || routes[0] == "" {
			igcinfo.RegisterTrack(w, r)
			return
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
