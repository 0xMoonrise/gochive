package handlers

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
)

func UploadFile(app *core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			slog.Error("failed to parse multipart form", "error", err)
			Error(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			slog.Error("something went wrong while uploading the a file", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong")
			return
		}
		defer file.Close()

		filename := filepath.Base(header.Filename)

		if header.Size == 0 {
			Error(w, http.StatusBadRequest, "file is empty")
			return
		}

		if !utils.ValidateFilename(filename) {
			Error(w, http.StatusBadRequest, "file type is not allowed")
			return
		}

		if utils.IsTooLong(filename) {
			Error(w, http.StatusBadRequest, "filename too long")
			return
		}

		tx, err := app.DB.Begin()
		if err != nil {
			slog.Error("failed to begin transaction", "error", err)
			Error(w, http.StatusBadRequest, "something went wrong")
			return
		}
		defer tx.Rollback()
		qtx := app.DB.Queries.WithTx(tx)

		id, err := qtx.InsertFile(r.Context(),
			database.InsertFileParams{
				Filename:  filename,
				Editorial: "Default",
			})
		if err != nil {
			slog.Error("error while trying to store metada file into the database", "error", err)
			Error(w, http.StatusBadRequest, "uploaded unsuccessful")
			return
		}

		sniff := make([]byte, 512)
		n, err := file.ReadAt(sniff, 0)
		if err != nil && err != io.EOF {
			slog.Error("error while trying to read 512 for content type detection", "error", err)
			Error(w, http.StatusBadRequest, "uploaded unsuccessful")
			return
		}
		contentType := utils.DetectContentType(sniff[:n])

		err = app.Storage.PutItem(
			r.Context(),
			path.Join("files", strconv.Itoa(id)),
			&core.Object{
				Length:      header.Size,
				ContentType: contentType,
				Reader:      file,
			},
		)
		if err != nil {
			slog.Error("error while trying to upload a file to storage", "error", err)
			Error(w, http.StatusBadRequest, "uploaded unsuccessful")
			return
		}

		if strings.HasSuffix(filename, ".md") {
			if err := tx.Commit(); err != nil {
				slog.Error("failed to commit transaction on markdown commit", "error", err)
				Error(w, http.StatusInternalServerError, "Something went wrong")
				return
			}
			JSON(w, http.StatusOK, SuccesUpload{
				Success: true,
				FileResponse: FileResponse{
					Id:        id,
					Filename:  filename,
					Editorial: "Default",
					Favorite:  false,
				},
			})
			return
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			slog.Error("failed to seek file back to start for thumbnail", "error", err)
			Error(w, http.StatusBadRequest, "uploaded unsuccessful")
			return
		}

		image := &bytes.Buffer{}
		if err := utils.MakeThumbnail(file, header.Size, 0, image); err != nil {
			slog.Error("error while trying to generate the thumbnail", "error", err)
			Error(w, http.StatusBadRequest, "uploaded unsuccessful")
			return
		}

		imageBytes := image.Bytes()
		imageReader := bytes.NewReader(imageBytes)

		objKey := path.Join("images", strconv.Itoa(id))
		err = app.Storage.PutItem(r.Context(), objKey, &core.Object{
			Length:      imageReader.Size(),
			ContentType: utils.DetectContentType(imageBytes),
			Reader:      io.NopCloser(imageReader),
		})
		if err != nil {
			slog.Error("error while trying to upload a image to storage", "error", err)
			Error(w, http.StatusBadRequest, "uploaded unsuccessful")
			return
		}

		if err := tx.Commit(); err != nil {
			slog.Error("failed to commit transaction at the end of upload handler", "error", err)
			Error(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		JSON(w, http.StatusOK, SuccesUpload{
			Success: true,
			FileResponse: FileResponse{
				Id:        id,
				Filename:  filename,
				Editorial: "Default",
				Favorite:  false,
			},
		})
	}
}
