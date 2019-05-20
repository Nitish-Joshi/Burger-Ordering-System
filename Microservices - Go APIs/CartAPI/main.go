package main

import (
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}

	server := NewServer()
	server.Run(":" + port)
}
