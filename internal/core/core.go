package core

import (
	"context"
	"database/sql"
	"io"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/database"
)

type Database struct {
	*sql.DB
	*database.Queries
}

type App struct {
	DB      Database
	Storage Store
	Config  *config.Config
}

type Object struct {
	Length      int64
	ContentType string
	Reader      io.ReadCloser
}

type Store interface {
	GetItem(ctx context.Context, objKey string) (obj *Object, err error)
	PutItem(ctx context.Context, objKey string, obj *Object) (err error)
	DelItem(ctx context.Context, objKey string) (err error)
}
