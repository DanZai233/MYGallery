package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mygallery/mygallery/internal/config"
)

// Storage 存储接口
type Storage interface {
	Upload(filename string, file io.Reader) (string, error)
	UploadThumbnail(filename string, file io.Reader) (string, error)
	Delete(path string) error
	GetURL(path string) string
}

var globalStorage Storage

// InitStorage 初始化存储
func InitStorage(cfg *config.Config) error {
	var err error
	
	switch cfg.Storage.Type {
	case "local":
		globalStorage, err = NewLocalStorage(cfg)
	case "s3":
		globalStorage, err = NewS3Storage(cfg)
	case "minio":
		globalStorage, err = NewMinIOStorage(cfg)
	case "aliyun":
		globalStorage, err = NewAliyunStorage(cfg)
	default:
		return fmt.Errorf("不支持的存储类型: %s", cfg.Storage.Type)
	}
	
	return err
}

// GetStorage 获取存储实例
func GetStorage() Storage {
	return globalStorage
}

// LocalStorage 本地存储实现
type LocalStorage struct {
	uploadDir    string
	thumbnailDir string
	urlPrefix    string
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(cfg *config.Config) (*LocalStorage, error) {
	// 创建上传目录
	if err := os.MkdirAll(cfg.Storage.Local.UploadDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(cfg.Storage.Local.ThumbnailDir, 0755); err != nil {
		return nil, err
	}
	
	return &LocalStorage{
		uploadDir:    cfg.Storage.Local.UploadDir,
		thumbnailDir: cfg.Storage.Local.ThumbnailDir,
		urlPrefix:    cfg.Storage.Local.URLPrefix,
	}, nil
}

// Upload 上传文件
func (s *LocalStorage) Upload(filename string, file io.Reader) (string, error) {
	filepath := filepath.Join(s.uploadDir, filename)
	
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	
	return filename, nil
}

// UploadThumbnail 上传缩略图到本地
func (s *LocalStorage) UploadThumbnail(filename string, file io.Reader) (string, error) {
	fp := filepath.Join(s.thumbnailDir, filename)
	dst, err := os.Create(fp)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}
	return filename, nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(path string) error {
	filepath := filepath.Join(s.uploadDir, path)
	return os.Remove(filepath)
}

// GetURL 获取文件URL
func (s *LocalStorage) GetURL(path string) string {
	return s.urlPrefix + "/" + path
}

// GetThumbnailDir 获取缩略图目录
func (s *LocalStorage) GetThumbnailDir() string {
	return s.thumbnailDir
}

// GetUploadDir 获取上传目录
func (s *LocalStorage) GetUploadDir() string {
	return s.uploadDir
}

