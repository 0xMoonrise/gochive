package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/utils"
	"github.com/spf13/cobra"
)

func newGenerateThumbnail() *cobra.Command {
	app := core.NewApp()
	return &cobra.Command{
		Use:   "generate [id]",
		Short: "Generate thumbnail(s)",
		Long:  "Provide an ID to generate a specific thumbnail, or omit the ID to generate all thumbnails.",
		Example: `  gochive generate 42
  gochive generate`,
		Args:                  cobra.RangeArgs(0, 1),
		DisableFlagsInUseLine: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.Run(
				core.StageConfig,
				core.StageDB,
				core.StageStorage,
			)
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return generateAllThumbnail(app)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			return generateOneThumbnail(ctx, app, args[0])
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return app.Cleanup()
		},
	}
}

func generateOneThumbnail(ctx context.Context, app *core.App, arg string) error {

	id, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}

	file, err := app.DB.Queries.GetArchive(ctx, id)
	if err != nil {
		return err
	}

	if filepath.Ext(file.Filename) == ".md" {
		slog.Warn("Cannot generate a thumbnail for a Markdown file", "id", file.ID)
		return nil
	}

	objKey := path.Join("files", arg)
	item, err := app.Storage.GetItem(ctx, objKey)
	if err != nil {
		return err
	}
	defer item.Reader.Close()

	fileData, err := io.ReadAll(item.Reader)
	if err != nil {
		return err
	}

	image := &bytes.Buffer{}
	if err := utils.MakeThumbnail(bytes.NewReader(fileData), item.Length, 0, image); err != nil {
		return err
	}

	imageBytes := image.Bytes()
	imageReader := bytes.NewReader(imageBytes)

	objKey = path.Join("images", arg)
	obj := &core.Object{
		Length:      imageReader.Size(),
		ContentType: http.DetectContentType(imageBytes[:512]),
		Reader:      io.NopCloser(imageReader),
	}

	if err := app.Storage.PutItem(ctx, objKey, obj); err != nil {
		return err
	}

	slog.Info("The thumbnial has been generate successfully", "id", arg)
	return nil
}

func generateAllThumbnail(app *core.App) error {
	const maxConcurrency = 8
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rows, err := app.DB.Queries.GetAllFiles(ctx)
	if err != nil {
		return err
	}

	sem := make(chan struct{}, maxConcurrency)
	for _, row := range rows {
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			id := strconv.Itoa(row.ID)
			if err := generateOneThumbnail(ctx, app, id); err != nil {
				slog.Error("An error ocurred while trying to generate a thumbnail", "error", err, "id", id)
			}
		}()
	}
	wg.Wait()
	return nil
}
