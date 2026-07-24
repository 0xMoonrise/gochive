package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const backupConcurrency = 8

func newBackupCmd() *cobra.Command {
	app := core.NewApp()
	var backupPath string

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Backup gochive data to the external backup location",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Run(
				core.StageConfig,
				core.StageDB,
				core.StageStorage,
			); err != nil {
				return err
			}
			if backupPath == "" {
				return errors.New("no backup path provided: use -p")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackup(cmd.Context(), app, backupPath)
		},
	}

	cmd.Flags().StringVarP(&backupPath, "path", "p", "", "backup destination path (overrides configured backup path)")
	return cmd
}

func newRestoreCmd() *cobra.Command {
	app := core.NewApp()
	var backupPath string

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore gochive data from the external backup location",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := app.Run(
				core.StageConfig,
				core.StageStorage,
			); err != nil {
				return err
			}
			if backupPath == "" {
				return errors.New("no backup path provided: use -p")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestore(cmd.Context(), app, backupPath)
		},
	}

	cmd.Flags().StringVarP(&backupPath, "path", "p", "", "backup source path (overrides configured backup path)")
	return cmd
}

func runBackup(ctx context.Context, app *core.App, backupPath string) error {
	for _, sub := range []string{"files", "images"} {
		if err := os.MkdirAll(filepath.Join(backupPath, sub), 0o755); err != nil {
			return fmt.Errorf("creating %s dir: %w", sub, err)
		}
	}

	files, err := app.DB.Queries.GetAllFiles(ctx)
	if err != nil {
		return fmt.Errorf("listing archive files: %w", err)
	}

	if err := backupObjects(ctx, app, backupPath, files); err != nil {
		return err
	}

	dbBackupPath := filepath.Join(backupPath, "gochive.db")
	if err := backupDatabase(ctx, app.DB.DB, dbBackupPath); err != nil {
		return fmt.Errorf("backing up database: %w", err)
	}

	slog.Info("backup completed", "path", backupPath, "files", len(files))
	return nil
}

func backupObjects(ctx context.Context, app *core.App, backupPath string, files []database.GetAllFilesRow) error {
	g := new(errgroup.Group)
	g.SetLimit(backupConcurrency)

	for _, file := range files {
		idStr := fmt.Sprintf("%d", file.ID)

		g.Go(func() error {
			return copyOne(ctx, app, backupPath, "files", idStr)
		})

		if !strings.HasSuffix(file.Filename, ".md") {
			g.Go(func() error {
				return copyOne(ctx, app, backupPath, "images", idStr)
			})
		}
	}
	return g.Wait()
}

func copyOne(ctx context.Context, app *core.App, backupPath, prefix, idStr string) error {
	objKey := path.Join(prefix, idStr)
	localPath := filepath.Join(backupPath, prefix, idStr)

	if _, err := os.Stat(localPath); err == nil {
		slog.Info("skipping, already backed up", "item", objKey)
		return nil
	}

	obj, err := app.Storage.GetItem(ctx, objKey)
	if err != nil {
		return fmt.Errorf("get %s: %w", objKey, err)
	}
	defer obj.Reader.Close()

	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("create %s: %w", localPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, obj.Reader); err != nil {
		return fmt.Errorf("copy %s: %w", objKey, err)
	}

	slog.Info("backup succeeded", "item", objKey)
	return nil
}

func backupDatabase(ctx context.Context, db *sql.DB, dest string) error {
	tmp := dest + ".tmp"
	if err := os.Remove(tmp); err != nil && !os.IsNotExist(err) {
		return err
	}

	if _, err := db.ExecContext(ctx, "VACUUM INTO ?", tmp); err != nil {
		return fmt.Errorf("vacuum into %s: %w", tmp, err)
	}

	return os.Rename(tmp, dest)
}

func runRestore(ctx context.Context, app *core.App, backupPath string) error {
	dbBackupPath := filepath.Join(backupPath, "gochive.db")
	if _, err := os.Stat(dbBackupPath); err != nil {
		return fmt.Errorf("backup database not found at %s: %w", dbBackupPath, err)
	}

	if err := restoreObjects(ctx, app, filepath.Join(backupPath, "files"), "files"); err != nil {
		return fmt.Errorf("restoring files: %w", err)
	}
	if err := restoreObjects(ctx, app, filepath.Join(backupPath, "images"), "images"); err != nil {
		return fmt.Errorf("restoring images: %w", err)
	}

	dbDestPath := filepath.Join(app.Config.Data, "gochive.db")
	if err := copyFile(dbBackupPath, dbDestPath); err != nil {
		return fmt.Errorf("copying database to %s: %w", dbDestPath, err)
	}

	slog.Info("restore completed", "database", dbDestPath)
	return nil
}

func restoreObjects(ctx context.Context, app *core.App, localDir, prefix string) error {
	entries, err := os.ReadDir(localDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	g := new(errgroup.Group)
	g.SetLimit(backupConcurrency)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		g.Go(func() error {
			objKey := path.Join(prefix, name)
			localPath := filepath.Join(localDir, name)

			f, err := os.Open(localPath)
			if err != nil {
				return fmt.Errorf("open %s: %w", localPath, err)
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				return fmt.Errorf("stat %s: %w", localPath, err)
			}

			if err := app.Storage.PutItem(ctx, objKey, &core.Object{
				Length: info.Size(),
				Reader: f,
			}); err != nil {
				return fmt.Errorf("put %s: %w", objKey, err)
			}

			slog.Info("restore succeeded", "item", objKey)
			return nil
		})
	}

	return g.Wait()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create %s: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %s to %s: %w", src, dst, err)
	}
	return out.Sync()
}
