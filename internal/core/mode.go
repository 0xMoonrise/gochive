package core

import (
	"errors"
	"log/slog"

	"github.com/0xMoonrise/gochive/internal/config"
)

func (app App) SetMode() (client Store, err error) {
	switch app.Config.Mode {
	case config.FS:
		client, err = app.NewfsClient()
		slog.Info("File system client has been selected")
		return
	case config.S3:
		client, err = app.NewS3Client()
		slog.Info("S3 client has been selected")
		return
	}
	return nil, errors.New("No mode was set")
}
