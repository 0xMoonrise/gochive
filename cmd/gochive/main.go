package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/server"
	_ "github.com/lib/pq"
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

	conn, err := config.LoadConfig()

	if err != nil {
		slog.Error("Database is not connected")
		log.Fatal(err)
	}

	defer conn.Close()
	db := database.New(conn)

	dumpImages("static/thumbnails/", db)

	server := server.NewServer(db)
	server.Run(addr)
}
