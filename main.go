package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {

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
