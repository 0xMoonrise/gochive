package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/0xMoonrise/gochive/internal/database"
)

func dumpImages(path string, db *database.Queries) {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	data, err := db.GetThumbnails(ctx)

	if err != nil {
		slog.Error("something went wrong while trying to dump the thumbnails")
	}

	defer cancel()

	for _, image := range data {

		path := filepath.Join(path, strconv.Itoa(int(image.ID)))

		_, err := os.Stat(path)

		if err == nil || !os.IsNotExist(err) {
			continue
		}

		os.WriteFile(path, image.ThumbnailImage, 0644)

	}
}
