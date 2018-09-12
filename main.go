package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/imt2681-assig1/igcinfo"
)

const rootPath = "/igcinfo/api/"

func handleRoutes(w http.ResponseWriter, r *http.Request) {
	// Get route
	route := strings.TrimPrefix(r.URL.Path, rootPath)

	if r.Method == http.MethodGet {
		// GET /api
		if route == "" {
			igcinfo.APIInfo(w, r)
			return
		}

		// GET /api/igc
		if match, _ := regexp.MatchString("^igc[/]?$", route); match {
			igcinfo.GetAllTracks(w, r)
			return
		}

		// GET /api/igc/{id}
		match := regexp.MustCompile("^igc/([0-9]+)[/]?$").FindStringSubmatch(route)
		if len(match) == 2 {
			id, _ := strconv.Atoi(match[1])
			igcinfo.GetTrackByID(w, r, id)
			return
		}

		// GET /api/igc/{id}/{field}
		match = regexp.MustCompile("^igc/([0-9]+)/([^/.]+)[/]?$").FindStringSubmatch(route)
		if len(match) == 3 {
			id, _ := strconv.Atoi(match[1])
			igcinfo.GetTrackField(w, r, id, match[2])
			return
		}
	}

	if r.Method == http.MethodPost {
		// POST /api/igc
		if match, _ := regexp.MatchString("^igc[/]?$", route); match {
			igcinfo.RegisterTrack(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT not set")
	}

	http.HandleFunc(rootPath, handleRoutes)

	// Start listen
	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
