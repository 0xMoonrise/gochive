package handlers

import (
	"net/http"

	"github.com/0xMoonrise/gochive/internal/core"
)

func Root(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, app.Templates, "index.html", http.StatusOK, map[string]any{
			"title": "Archive",
		})
	}
}
