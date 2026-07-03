package main

import (
	"errors"
	"log/slog"

	"github.com/0xMoonrise/gochive/internal/core"
)

const (
	fileSystem = iota + 1
	S3
)

func setMode(mode int) (client core.Store, err error) {
	switch mode {
	case fileSystem:
		client, err = core.NewfsClient()
		slog.Info("File system client was set")
		return
	case S3:
		client, err = core.NewS3Client()
		slog.Info("S3 client was set")
		return
	}
	return nil, errors.New("No mode was set")
}
