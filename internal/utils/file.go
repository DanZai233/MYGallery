package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// GenerateFilename 生成唯一文件名
func GenerateFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	hash := md5.Sum([]byte(fmt.Sprintf("%s%d", originalName, timestamp)))
	return fmt.Sprintf("%x%s", hash, ext)
}

// IsAllowedImageType 检查是否是允许的图片类型
func IsAllowedImageType(mimeType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if strings.EqualFold(mimeType, allowed) {
			return true
		}
	}
	return false
}

// GetFileSize 获取上传文件大小
func GetFileSize(file multipart.File) (int64, error) {
	// 移动到文件末尾获取大小
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	// 移回文件开头
	_, err = file.Seek(0, io.SeekStart)
	return size, err
}

