package storage

import (
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mygallery/mygallery/internal/config"
)

// S3Storage AWS S3 存储实现
type S3Storage struct {
	client    *s3.S3
	bucket    string
	urlPrefix string
}

// NewS3Storage 创建 S3 存储
func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	awsConfig := &aws.Config{
		Region:      aws.String(cfg.Storage.S3.Region),
		Credentials: credentials.NewStaticCredentials(cfg.Storage.S3.AccessKey, cfg.Storage.S3.SecretKey, ""),
	}
	
	if cfg.Storage.S3.Endpoint != "" {
		awsConfig.Endpoint = aws.String(cfg.Storage.S3.Endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}
	
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("创建 S3 会话失败: %w", err)
	}
	
	return &S3Storage{
		client:    s3.New(sess),
		bucket:    cfg.Storage.S3.Bucket,
		urlPrefix: cfg.Storage.S3.URLPrefix,
	}, nil
}

// Upload 上传文件到 S3
func (s *S3Storage) Upload(filename string, file io.Reader) (string, error) {
	// 读取文件内容
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}
	
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),
	})
	
	if err != nil {
		return "", fmt.Errorf("上传到 S3 失败: %w", err)
	}
	
	return filename, nil
}

// Delete 从 S3 删除文件
func (s *S3Storage) Delete(path string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	return err
}

// GetURL 获取文件 URL
func (s *S3Storage) GetURL(path string) string {
	return fmt.Sprintf("%s/%s", s.urlPrefix, path)
}

