package app

import (
	"context"
	"io"

	"github.com/0xMoonrise/gochive/internal/database"
)

type App struct {
	Db      *database.Queries
	Storage Store
}

type Store interface {
	GetItem(ctx context.Context, objKey string) (length int64, contentType string, reader io.ReadCloser, err error)
	PutItem(ctx context.Context, objKey string, length int64, contentType string, reader io.ReadCloser) (err error)
}
