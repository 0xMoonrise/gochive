package core

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/0xMoonrise/gochive/internal/utils"
)

type fsClient struct {
	Path string
}

func (c *fsClient) GetItem(ctx context.Context, objKey string) (obj *Object, err error) {
	pathTo := path.Join(c.Path, objKey)

	file, err := os.OpenFile(pathTo, os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	info, err := file.Stat()
	if err != nil {
		return
	}

	obj = &Object{}
	size := info.Size()
	obj.Length = size
	sniff := min(size, 512)
	buffer := make([]byte, sniff)
	n, err := file.ReadAt(buffer, 0)
	if err != nil && err != io.EOF {
		return
	}

	obj.ContentType = utils.DetectContentType(buffer[:n])

	obj.Reader = file
	return
}

func (c *fsClient) PutItem(ctx context.Context, objKey string, obj *Object) (err error) {
	path := path.Join(c.Path, objKey)

	f, err := os.Create(path)
	if err != nil {
		return
	}

	defer f.Close()
	_, err = io.Copy(f, obj.Reader)

	return
}

func (c *fsClient) DelItem(
	ctx context.Context,
	objKey string,
) (err error) {

	path := path.Join(c.Path, objKey)
	if err = os.Remove(path); err != nil {
		return
	}
	return nil
}

func (app App) NewfsClient() (client *fsClient, err error) {

	if err := os.MkdirAll(app.Config.FS.Root, 0755); err != nil {
		return nil, err
	}

	client = &fsClient{
		Path: app.Config.FS.Root,
	}

	dirs := []string{"images", "files"}
	for _, dir := range dirs {
		path := filepath.Join(client.Path, dir)
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
