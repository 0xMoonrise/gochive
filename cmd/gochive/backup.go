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
	"sync"

	"strings"

	"github.com/0xMoonrise/gochive/internal/core"
	"github.com/0xMoonrise/gochive/internal/database"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
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
				core.StageStorage); err != nil {
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
			if err := app.Run(core.StageConfig,
				core.StageStorage); err != nil {
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
	if err := os.MkdirAll(filepath.Join(backupPath, "files"), 0o755); err != nil {
		return fmt.Errorf("creating files backup dir: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(backupPath, "images"), 0o755); err != nil {
		return fmt.Errorf("creating images backup dir: %w", err)
	}

	files, err := app.DB.Queries.GetAllFiles(ctx)
	if err != nil {
		return fmt.Errorf("listing archive files: %w", err)
	}

	if err := backupObjects(ctx, app, backupPath, files); err != nil {
		return err
	}

	dbBackupPath := filepath.Join(backupPath, "gochive.db")
	if err := backupDatabase(ctx, app.DB.DB, dbBackupPath, len(files)); err != nil {
		return fmt.Errorf("backing up database: %w", err)
	}

	configSrc := filepath.Join(app.Config.Data, "config.toml")
	configDst := filepath.Join(backupPath, "config.toml")
	if err := copyFile(configSrc, configDst); err != nil {
		return fmt.Errorf("backing up config.toml: %w", err)
	}

	slog.Info("backup completed", "path", backupPath, "files", len(files))
	return nil
}

func backupObjects(ctx context.Context, app *core.App, backupPath string, files []database.GetAllFilesRow) error {
	sem := make(chan struct{}, backupConcurrency)
	var wg sync.WaitGroup
	errCh := make(chan error, len(files)*2)

	copyOne := func(objKey, localPath string) {
		defer wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()

		obj, err := app.Storage.GetItem(ctx, objKey)
		if err != nil {
			errCh <- fmt.Errorf("get %s: %w", objKey, err)
			return
		}
		defer obj.Reader.Close()

		f, err := os.Create(localPath)
		if err != nil {
			errCh <- fmt.Errorf("create %s: %w", localPath, err)
			return
		}
		defer f.Close()

		if _, err := io.Copy(f, obj.Reader); err != nil {
			errCh <- fmt.Errorf("copy %s: %w", objKey, err)
		}
		slog.Info("An item has been backup", "id", objKey)
	}

	for _, file := range files {
		idStr := fmt.Sprintf("%d", file.ID)
		hasImage := !strings.HasSuffix(file.Filename, ".md")

		wg.Add(1)
		go copyOne(path.Join("files", idStr), filepath.Join(backupPath, "files", idStr))

		if hasImage {
			wg.Add(1)
			go copyOne(path.Join("images", idStr), filepath.Join(backupPath, "images", idStr))
		}
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func backupDatabase(ctx context.Context, db *sql.DB, dest string, newCount int) error {
	if existingCount, err := countArchiveRows(dest); err == nil {
		if newCount < existingCount {
			return fmt.Errorf(
				"refusing to replace existing backup db (%d rows) with a smaller one (%d rows) at %s: "+
					"the source database may be incomplete or pointing to the wrong path",
				existingCount, newCount, dest,
			)
		}
	}

	tmp := dest + ".tmp"
	_ = os.Remove(tmp)

	if _, err := db.ExecContext(ctx, "VACUUM INTO ?", tmp); err != nil {
		return fmt.Errorf("vacuum into %s: %w", tmp, err)
	}

	return os.Rename(tmp, dest)
}

func countArchiveRows(dest string) (int, error) {
	if _, err := os.Stat(dest); err != nil {
		return 0, err
	}

	conn, err := sql.Open("sqlite3", "file:"+dest+"?mode=ro")
	if err != nil {
		return 0, fmt.Errorf("opening existing backup for row count: %w", err)
	}
	defer conn.Close()

	var count int
	if err := conn.QueryRow("SELECT COUNT(*) FROM archive").Scan(&count); err != nil {
		return 0, fmt.Errorf("counting rows in existing backup: %w", err)
	}
	return count, nil
}

func runRestore(ctx context.Context, app *core.App, backupPath string) error {
	filesDir := filepath.Join(backupPath, "files")
	imagesDir := filepath.Join(backupPath, "images")
	dbBackupPath := filepath.Join(backupPath, "gochive.db")

	if _, err := os.Stat(dbBackupPath); err != nil {
		return fmt.Errorf("backup database not found at %s: %w", dbBackupPath, err)
	}

	if err := restoreObjects(ctx, app, filesDir, "files"); err != nil {
		return fmt.Errorf("restoring files: %w", err)
	}
	if err := restoreObjects(ctx, app, imagesDir, "images"); err != nil {
		return fmt.Errorf("restoring images: %w", err)
	}

	dbDestPath := filepath.Join(app.Config.Data, "gochive.db")
	if err := copyFile(dbBackupPath, dbDestPath); err != nil {
		return fmt.Errorf("copying database to %s: %w", dbDestPath, err)
	}

	configSrc := filepath.Join(backupPath, "config.toml")
	configDst := filepath.Join(app.Config.Data, "config.toml")
	if err := copyFile(configSrc, configDst); err != nil {
		return fmt.Errorf("restoring config.toml: %w", err)
	}

	slog.Info("restore completed", "database", dbDestPath, "config", configDst)
	return nil
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

func restoreObjects(ctx context.Context, app *core.App, localDir, prefix string) error {
	entries, err := os.ReadDir(localDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	sem := make(chan struct{}, backupConcurrency)
	var wg sync.WaitGroup
	errCh := make(chan error, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			localPath := filepath.Join(localDir, name)
			f, err := os.Open(localPath)
			if err != nil {
				errCh <- fmt.Errorf("open %s: %w", localPath, err)
				return
			}
			defer f.Close()

			info, err := f.Stat()
			if err != nil {
				errCh <- fmt.Errorf("stat %s: %w", localPath, err)
				return
			}

			objKey := path.Join(prefix, name)
			err = app.Storage.PutItem(ctx, objKey, &core.Object{
				Length: info.Size(),
				Reader: f,
			})
			if err != nil {
				errCh <- fmt.Errorf("put %s: %w", objKey, err)
			}
			slog.Info("An item has been restored", "id", name)
		}(entry.Name())
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
