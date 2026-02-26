package storage

import "mime/multipart"

// Storage 定义了文件存储的通用接口
type Storage interface {
	// Save 保存文件
	// subDir: 子目录，例如 "character/photos" 或 "user/photos"
	// 返回值: 文件的相对路径或云端 URL，以及 error
	Save(file *multipart.FileHeader, subDir string, fileNamePrefix string) (string, error)

	// Delete 删除文件
	// fileKey: 数据库中存储的相对路径或 Object Key
	Delete(fileKey string) error
}