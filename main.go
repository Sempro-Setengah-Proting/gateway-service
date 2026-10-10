package main

import (
	"log"

	"gatewayservice/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
