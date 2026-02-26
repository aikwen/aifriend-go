package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

// LocalStorage 基于本地文件系统的存储实现
type LocalStorage struct {
	baseDir string // 根目录，例如 "media"
}

// NewLocalStorage 实例化本地存储
func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{
		baseDir: baseDir,
	}
}

// Save 将文件存储到 baseDir/subDir 目录下
// 并返回 文件相对 baseDir 的路径
func (l *LocalStorage) Save(file *multipart.FileHeader, subDir string, fileNamePrefix string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	
	var fileName string
	if fileNamePrefix != "" {
		fileName = fileNamePrefix + "_" + uuid.New().String() + ext
	} else {
		fileName = uuid.New().String() + ext
	}

	uploadDir := filepath.Join(l.baseDir, subDir)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}
	fullDiskPath := filepath.Join(uploadDir, fileName)

	dst, err := os.Create(fullDiskPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return path.Join(subDir, fileName), nil
}

func (l *LocalStorage) Delete(fileKey string) error {
	if fileKey == "" {
		return nil
	}

	diskPath := filepath.Join(l.baseDir, fileKey)
	return os.Remove(diskPath)
}