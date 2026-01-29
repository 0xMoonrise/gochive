package utils

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/0xMoonrise/gochive/internal/config"
)

func DumpThumbnails(c *config.Conf) {

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	thumbnails, err := c.Db.GetThumbnails(ctx)

	if err != nil {
		slog.Error("failed to fetch thumbnails", "err", err)
		return
	}

	if _, err := os.Stat(config.THUMB_PATH); os.IsNotExist(err) {
		slog.Warn("thumbnail directory does not exist, creating", "path", config.THUMB_PATH)

		if err := os.MkdirAll(config.THUMB_PATH, 0755); err != nil {
			slog.Error("failed to create thumbnail directory", "path", config.THUMB_PATH, "err", err)
			return
		}
	}

	// goroutines for parallelism?

	for _, thumbnail := range thumbnails {

		thumbToWrite := filepath.Join(
			config.THUMB_PATH,
			strconv.Itoa(int(thumbnail.ID)),
		)

		if err := os.WriteFile(thumbToWrite, thumbnail.ThumbnailImage, 0644); err != nil {
			slog.Debug("thumbnail write skipped", "id", thumbnail.ID)
		}

	}
}
