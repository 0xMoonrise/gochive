package main

import (
	"log/slog"
	"net"
	"os"

	c "github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/server"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func run() error {

	if err := godotenv.Load(); err != nil {
		slog.Info(".env not loaded")
	}

	conf := &c.Conf{}
	db, err := connectDB()

	if err != nil {
		return err
	}

	defer db.Close()
	database := database.New(db)

	if err := bootSchema(database); err != nil {
		return err
	}

	if err := bootMigrations(db); err != nil {
		return err
	}

	conf.Db = database
	conf.S3 = os.Getenv("S3_URL")

	server := server.NewServer(conf)
	addr := net.JoinHostPort(os.Getenv("HOST"), os.Getenv("PORT"))

	if err := server.Run(addr); err != nil {
		return err
	}

	return nil

}

func main() {

	if err := run(); err != nil {
		slog.Error("fatal:", "error", err)
		os.Exit(1)
	}

}
