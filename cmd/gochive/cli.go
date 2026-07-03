package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
)

func cli(app *core.App) error {
	if len(os.Args) < 2 {
		return nil
	}

	cmd := os.Args[1]

	switch cmd {
	case "normalize":
		runNormalize(app)
		return errors.New("End process")
	case "server":
		return nil
	}

	return errors.New("Uknwon command")
}

func runNormalize(app *core.App) {
	rows, err := app.Db.GetAllFiles(context.Background())
	if err != nil {
		log.Fatal("Something went wrong while trying to connect to db")
	}
	for _, file := range rows {
		key := path.Join("files", file.Filename)

		f, err := app.Storage.GetItem(context.Background(), key)
		if err != nil {
			log.Fatal(err)
		}

		id := strconv.Itoa(file.ID)
		if err := writeFromReader(f.Reader, path.Join(config.ROOT, "files", id)); err != nil {
			log.Fatal(err)
		}

		if err := app.Storage.DelItem(context.Background(), key); err != nil {
			log.Fatal(err)
		}

	}
}

func writeFromReader(r io.Reader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	return nil
}
