package main

import (
	"os"

	"github.com/haakonleg/imt2681-assig1/igcinfo"
)

const defaultPort = "8080"

func main() {
	// Get port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}

	// Configure and start the API
	app := igcinfo.App{
		ListenPort: port}
	app.StartServer()
}
