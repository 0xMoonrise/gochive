package app

import (
	"github.com/0xMoonrise/gochive/internal/database"
)

type App struct {
	Db       *database.Queries
	S3Client *Client
}
