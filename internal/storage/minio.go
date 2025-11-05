package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mygallery/mygallery/internal/config"
)

// MinIOStorage MinIO 存储实现
type MinIOStorage struct {
	client    *minio.Client
	bucket    string
	urlPrefix string
}

// NewMinIOStorage 创建 MinIO 存储
func NewMinIOStorage(cfg *config.Config) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Storage.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Storage.MinIO.AccessKey, cfg.Storage.MinIO.SecretKey, ""),
		Secure: cfg.Storage.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 MinIO 客户端失败: %w", err)
	}
	
	// 确保 bucket 存在
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Storage.MinIO.Bucket)
	if err != nil {
		return nil, fmt.Errorf("检查 bucket 失败: %w", err)
	}
	
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Storage.MinIO.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("创建 bucket 失败: %w", err)
		}
	}
	
	return &MinIOStorage{
		client:    client,
		bucket:    cfg.Storage.MinIO.Bucket,
		urlPrefix: cfg.Storage.MinIO.URLPrefix,
	}, nil
}

// Upload 上传文件到 MinIO
func (s *MinIOStorage) Upload(filename string, file io.Reader) (string, error) {
	ctx := context.Background()
	
	_, err := s.client.PutObject(ctx, s.bucket, filename, file, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("上传到 MinIO 失败: %w", err)
	}
	
	return filename, nil
}

// Delete 从 MinIO 删除文件
func (s *MinIOStorage) Delete(path string) error {
	ctx := context.Background()
	return s.client.RemoveObject(ctx, s.bucket, path, minio.RemoveObjectOptions{})
}

// GetURL 获取文件 URL
func (s *MinIOStorage) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", s.urlPrefix, path)
}

