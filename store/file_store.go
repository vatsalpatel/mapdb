package store

import (
	"os"
	"sync"
)

type FileStore struct {
	mu       sync.Mutex
	fileName string
	file     *os.File
}

func NewFileStore(fileName string) *FileStore {
	return &FileStore{
		fileName: fileName,
	}
}

func (fs *FileStore) Append(data []byte) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	var err error
	if fs.file == nil {
		fs.file, err = os.OpenFile(fs.fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			return err
		}
	}

	data = append(data, "\r\n"...)
	if _, err = fs.file.Write(data); err != nil {
		return err
	}
	if err = fs.file.Sync(); err != nil {
		return err
	}

	return nil
}

func (fs *FileStore) ReadAll() ([]byte, error) {
	if fs.file != nil {
		fs.file.Close()
		fs.file = nil
	}
	return os.ReadFile(fs.fileName)
}

func (fs *FileStore) Clear() error {
	var err error
	if fs.file == nil {
		fs.file, err = os.OpenFile(fs.fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
		if err != nil {
			return err
		}
	}
	return fs.file.Truncate(0)
}
