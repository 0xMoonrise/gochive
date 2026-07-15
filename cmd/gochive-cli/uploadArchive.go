package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/0xMoonrise/gochive/internal/config"
	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/spf13/cobra"
)

func newUploadArchive() *cobra.Command {
	app := core.NewApp()
	return &cobra.Command{
		Use:   "upload_archive <url>",
		Short: "Upload an archive from a URL",
		Long:  "Downloads an archive from the specified URL and uploads it to the archive.",
		Example: `gochive upload_archive https://example.com/archive.pdf
gochive upload_archive https://example.com/archive.md`,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(
				core.StageConfig,
				core.StageDB,
				core.StageStorage,
			)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := isValidURL(args[0])

			if err != nil {
				return err
			}

			return downloadFile(url, app)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return app.Cleanup()
		},
	}
}

func isValidURL(rawURL string) (*url.URL, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("scheme not supported")
	}

	if u.Host == "" {
		return nil, errors.New("url has no host")
	}

	if !utils.ValidateFilename(path.Ext(u.Path)) {
		return nil, errors.New("extension not allowed")
	}

	return u, nil
}

type progressReader struct {
	reader     io.Reader
	total      int64
	read       int64
	lastReport int
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.read += int64(n)

	if pr.total > 0 {
		percent := int(float64(pr.read) / float64(pr.total) * 100)
		if percent != pr.lastReport {
			fmt.Printf("\rDownload progress: %d%%", percent)
			pr.lastReport = percent
		}
	}

	return n, err
}

func downloadFile(url *url.URL, app *core.App) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("an error has occurred: " + res.Status)
	}

	pr := &progressReader{reader: res.Body, total: res.ContentLength}
	limitedReader := io.LimitReader(pr, config.MAX_UPLOAD_SIZE+1)
	buffer := bytes.Buffer{}

	bufferLength, err := io.Copy(&buffer, limitedReader)
	if err != nil {
		return err
	}

	file := buffer.Bytes()
	fmt.Println()
	if bufferLength > config.MAX_UPLOAD_SIZE {
		return errors.New("file exceeds maximum allowed size of 60MB")
	}

	tx, err := app.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	filename := path.Base(url.Path)
	slog.Info("Download success", "filename", filename)

	qtx := app.DB.Queries.WithTx(tx)
	id, err := qtx.InsertFile(ctx, database.InsertFileParams{
		Filename:  filename,
		Editorial: "Default",
	})

	if err != nil {
		return err
	}

	key := strconv.Itoa(id)
	objKey := path.Join("files", key)
	obj := &core.Object{
		Length:      bufferLength,
		ContentType: utils.DetectContentType(file),
		Reader:      io.NopCloser(bytes.NewReader(file)),
	}

	if err := app.Storage.PutItem(ctx, objKey, obj); err != nil {
		return err
	}
	if path.Ext(filename) == ".md" {
		slog.Info("A Markdown file does not require a thumbnail to be generated.")
		slog.Info("File successfully uploaded", "id", key)
		return tx.Commit()
	}

	image := &bytes.Buffer{}
	if err := utils.MakeThumbnail(bytes.NewReader(file), bufferLength, 0, image); err != nil {
		return err
	}

	imageBytes := image.Bytes()
	objKey = path.Join("images", key)
	obj = &core.Object{
		Length:      int64(image.Len()),
		ContentType: utils.DetectContentType(imageBytes),
		Reader:      io.NopCloser(bytes.NewReader(imageBytes)),
	}

	if err := app.Storage.PutItem(ctx, objKey, obj); err != nil {
		return err
	}

	slog.Info("File successfully uploaded", "id", key)
	return tx.Commit()
}
