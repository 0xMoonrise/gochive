package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
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

	const maxWorkers = 8
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	for _, file := range rows {
		wg.Add(1)
		go func(database.GetAllFilesRow) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := normalizeFile(app, file); err != nil {
				slog.Error("failed to normalize file", "filename", file.Filename, "err", err)
				return
			}
			slog.Info("file successfully normalized", "filename", file.Filename)
		}(file)
	}

	wg.Wait()
}

func normalizeFile(app *core.App, file database.GetAllFilesRow) error {
	key := path.Join("files", file.Filename)

	f, err := app.Storage.GetItem(context.Background(), key)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	defer f.Reader.Close()

	id := strconv.Itoa(file.ID)
	if err := writeFromReader(f.Reader, path.Join(config.ROOT, "files", id)); err != nil {
		return fmt.Errorf("write from reader: %w", err)
	}

	if err := app.Storage.DelItem(context.Background(), key); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
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
