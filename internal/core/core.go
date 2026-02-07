package core

import (
	"context"
	"io"

	"github.com/0xMoonrise/gochive/internal/database"
)

type App struct {
	Db      *database.Queries
	Storage Store
}

type Object struct {
	Length      int64
	ContentType string
	Reader      io.ReadCloser
}

type Store interface {
	GetItem(ctx context.Context, objKey string) (obj *Object, err error)
	PutItem(ctx context.Context, objKey string, obj *Object) (err error)
}
