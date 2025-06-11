package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/0xMoonrise/gochive/internal/server"
)

func main() {

	slog.Info("Init app")
	err := Init() // Should I just load from file or also env vars?

	if err != nil {
		log.Fatal("An error occured:", err)
	}
	
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	addr := fmt.Sprintf("%s:%s", host, port)
	
	server := server.NewServer()
	server.Run(addr)
}



