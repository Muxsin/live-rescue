package repositories

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type MediaRepository struct {
}

func NewMedia() *MediaRepository {
	return &MediaRepository{}
}

func (r *MediaRepository) StoreFile(uf multipart.File, ufh *multipart.FileHeader, filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, uf); err != nil {
		os.Remove(filePath)
		return err
	}

	return nil
}

func (r *MediaRepository) DeleteFile(filePath string) error {
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
