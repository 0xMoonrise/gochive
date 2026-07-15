package core

import (
	"errors"

	"github.com/0xMoonrise/gochive/internal/config"
)

func (app App) SetMode() (client Store, err error) {
	switch app.Config.Mode {
	case config.FS:
		client, err = app.NewfsClient()
		return
	case config.S3:
		client, err = app.NewS3Client()
		return
	}
	return nil, errors.New("No mode was set")
}

func (app App) ModeToString() string {
	switch app.Config.Mode {
	case config.FS:
		return "File System"
	case config.S3:
		return "S3"
	}
	return ""
}
