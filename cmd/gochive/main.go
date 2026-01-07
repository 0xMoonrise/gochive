package main

import (
  "database/sql"
  "os"
	"log/slog"

	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/server"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const MAX_RETRIES = 3

var (
  conn *sql.DB
  err error
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		slog.Info(".env not loaded")
	}

	conn, err := InitDataBase()
	if err != nil {
		panic(err)
	}
	
	defer conn.Close()
	database := database.New(conn)

	if err := loadDB(conn, database); err != nil {
		panic(err)
	}

	utils.DumpThumbnails(database)

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	server := server.NewServer(database)
	server.Run(host + ":" + port)

}

