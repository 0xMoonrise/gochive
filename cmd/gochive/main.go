package main

import (
	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/server"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	conn, err := config.InitDataBase()

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	database := database.New(conn)
	dumpThumbnails(database)

	server := server.NewServer(database)
	server.Run(host + ":" + port)

}
