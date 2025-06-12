package handlers

import (
	"github.com/0xMoonrise/gochive/internal/database"
)

type DBhdlr struct {
	Query *database.Queries
}

func Handler(db *database.Queries) *DBhdlr {
	return &DBhdlr{Query: db}
}

