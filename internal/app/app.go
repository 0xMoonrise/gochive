package app

import (
	"io"

	"github.com/0xMoonrise/gochive/internal/database"
)

type App struct {
	Db      *database.Queries
	Storage *Client
}

type Store interface {
	GetItem(name string) (length int64, contentType string, reader io.Reader)
}
