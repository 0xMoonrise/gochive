package config

import "github.com/0xMoonrise/gochive/internal/database"

type Conf struct {
	Db *database.Queries
	S3 string
}
