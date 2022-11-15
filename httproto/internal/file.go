package internal

import (
	"bytes"
	"errors"
	"io/fs"
	"io/ioutil"
)

var (
	errFoundAll = errors.New("found all")
	ErrReadOnly = errors.New("readonly fs")
)

type WebDAVFile struct {
	FileSystem fs.FS
	File       fs.File
	FileReader *bytes.Reader
}

func NewWebDAVFile(fileSystem fs.FS, file fs.File) (*WebDAVFile, error) {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &WebDAVFile{FileSystem: fileSystem, File: file, FileReader: bytes.NewReader(b)}, nil
}

func (wf *WebDAVFile) Close() error {
	return wf.File.Close()
}

func (wf *WebDAVFile) Read(p []byte) (n int, err error) {
	return wf.FileReader.Read(p)
}

func (wf *WebDAVFile) Readdir(count int) ([]fs.FileInfo, error) {
	info := []fs.FileInfo{}
	found := 0

	err := fs.WalkDir(wf.FileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if found >= count {
			return errFoundAll
		}

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		f, err := wf.FileSystem.Open(path)
		if err != nil {
			return err
		}

		fileInfo, err := f.Stat()
		if err != nil {
			return err
		}

		info = append(info, fileInfo)
		found++

		return nil
	})

	if err != errFoundAll && err != nil {
		return nil, err
	}

	return info, nil
}

func (wf *WebDAVFile) Seek(offset int64, whence int) (int64, error) {
	return wf.FileReader.Seek(offset, whence)
}

func (wf *WebDAVFile) Stat() (fs.FileInfo, error) {
	return wf.File.Stat()
}

func (wf *WebDAVFile) Write(p []byte) (n int, err error) {
	return 0, ErrReadOnly
}
