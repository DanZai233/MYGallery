package storage

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mygallery/mygallery/internal/config"
)

// AliyunStorage 阿里云 OSS 存储实现
type AliyunStorage struct {
	bucket    *oss.Bucket
	urlPrefix string
}

// NewAliyunStorage 创建阿里云存储
func NewAliyunStorage(cfg *config.Config) (*AliyunStorage, error) {
	client, err := oss.New(
		cfg.Storage.Aliyun.Endpoint,
		cfg.Storage.Aliyun.AccessKey,
		cfg.Storage.Aliyun.SecretKey,
	)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云 OSS 客户端失败: %w", err)
	}
	
	bucket, err := client.Bucket(cfg.Storage.Aliyun.Bucket)
	if err != nil {
		return nil, fmt.Errorf("获取 bucket 失败: %w", err)
	}
	
	return &AliyunStorage{
		bucket:    bucket,
		urlPrefix: cfg.Storage.Aliyun.URLPrefix,
	}, nil
}

// Upload 上传文件到阿里云 OSS
func (s *AliyunStorage) Upload(filename string, file io.Reader) (string, error) {
	// 读取文件内容
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}
	
	err := s.bucket.PutObject(filename, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return "", fmt.Errorf("上传到阿里云 OSS 失败: %w", err)
	}
	
	return filename, nil
}

// UploadThumbnail 上传缩略图到阿里云 OSS
func (s *AliyunStorage) UploadThumbnail(filename string, file io.Reader) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}
	key := "thumbnails/" + filename
	err := s.bucket.PutObject(key, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return "", fmt.Errorf("上传缩略图到阿里云 OSS 失败: %w", err)
	}
	return filename, nil
}

// Delete 从阿里云 OSS 删除文件
func (s *AliyunStorage) Delete(path string) error {
	return s.bucket.DeleteObject(path)
}

// GetURL 获取文件 URL
func (s *AliyunStorage) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", s.urlPrefix, path)
}

