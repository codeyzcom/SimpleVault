package storage

import (
	"log/slog"
	"os"
	"path/filepath"
)

type FileStorage struct {
	base string
}

func NewFileStorage(base string) *FileStorage {
	err := os.MkdirAll(base, 0700)
	if err != nil {
		slog.Error("NewFileStorage create with error", "detail", err.Error())
	}
	return &FileStorage{base: base}
}

func (fs *FileStorage) Read(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(fs.base, name))
}

func (fs *FileStorage) Write(name string, data []byte) error {
	tmp := filepath.Join(fs.base, name+".tmp")
	final := filepath.Join(fs.base, name)

	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, final)
}
