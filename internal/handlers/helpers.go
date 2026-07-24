package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0xMoonrise/gochive/internal/core"
)

func render(w http.ResponseWriter, tmpl *template.Template, name string, status int, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		slog.Error("template render failed", "template", name, "error", err)
	}
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("failed to encode json response", "error", err)
	}
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"status": msg})
}

func DecodeJSON[T any](w http.ResponseWriter, r *http.Request, maxBytes int64) (T, error) {
	var v T
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&v)
	return v, err
}

func fromStorageObject(w http.ResponseWriter, obj *core.Object) error {
	w.Header().Set("Content-Type", obj.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(obj.Length, 10))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, obj.Reader); err != nil {
		return err
	}

	return nil
}
