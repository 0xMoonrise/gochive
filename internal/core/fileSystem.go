package core

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/0xMoonrise/gochive/internal/utils"
)

type fsClient struct {
	Path string
}

func (c *fsClient) GetItem(ctx context.Context, objKey string) (
	obj *Object,
	err error,
) {

	pathTo := path.Join(c.Path, objKey)
	file, err := os.OpenFile(pathTo, os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			file.Close()
		}
	}()

	obj = &Object{}
	info := utils.Must(file.Stat())

	obj.Length = info.Size()
	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	obj.ContentType = http.DetectContentType(buffer[:n])
	if obj.ContentType == "" {
		obj.ContentType = "application/octet-stream"
	}
	obj.Reader = file

	return
}

func (c *fsClient) PutItem(
	ctx context.Context,
	objKey string,
	obj *Object,
) (
	// return
	err error,
) {
	pathTo := path.Join(c.Path, objKey)

	f, err := os.Create(pathTo)
	if err != nil {
		return
	}

	defer f.Close()
	_, err = io.Copy(f, obj.Reader)

	return
}

func NewfsClient() (client *fsClient, err error) {

	dbDir := "/opt/gochive"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	client = &fsClient{
		Path: "/opt/gochive/",
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
