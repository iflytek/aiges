package internal

import (
	"context"
	"io/fs"
	"os"

	"golang.org/x/net/webdav"
)

type WebDAVFileSystem struct {
	FileSystem fs.FS
}

func NewWebDAVFileSystemFromFS(fileSystem fs.FS) *WebDAVFileSystem {
	return &WebDAVFileSystem{FileSystem: fileSystem}
}

func (wfs *WebDAVFileSystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return ErrReadOnly
}

func (wfs *WebDAVFileSystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	file, err := wfs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return NewWebDAVFile(wfs.FileSystem, file)
}

func (wfs *WebDAVFileSystem) RemoveAll(ctx context.Context, name string) error {
	return ErrReadOnly
}

func (wfs *WebDAVFileSystem) Rename(ctx context.Context, oldName, newName string) error {
	return ErrReadOnly
}

func (wfs *WebDAVFileSystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	file, err := wfs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return file.Stat()
}
