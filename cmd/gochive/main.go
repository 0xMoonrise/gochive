package main

import (
	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/server"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	
	conn, err := config.InitDataBase()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	database := database.New(conn)

	if err := dbManager(conn, database); err != nil {
		log.Fatal(err)
	}
	
	utils.DumpThumbnails(database)
	
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	
	server := server.NewServer(database)
	server.Run(host + ":" + port)

}
