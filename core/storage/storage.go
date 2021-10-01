package storage

import (
	"io"
	"time"
)

type Driver interface {
	GetContent(path string) ([]byte, error)
	PutContent(path string, content []byte) error
	Reader(path string, offset int64) (io.ReadCloser, error)
	Writer(path string, append bool) (io.WriteCloser, error)
	Stat(path string) (FileInfo, error)
	List(path string) ([]string, error)
	Move(sourcePath, destPath string) error
	Delete(path string) error
	Walk(path string, f WalkFn) error
}

type Factory interface {
	Create(params map[string]interface{}) (Driver, error)
}

type FileInfo interface {
	Name() string
	Size() int64
	ModTime() time.Time
	IsDir() bool
}

type WalkFn func(fileInfo FileInfo) error
