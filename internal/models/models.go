package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"" json:"email"`
	Role      string         `gorm:"default:admin" json:"role"`
}

// Photo 照片模型
type Photo struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 基本信息
	Filename     string `gorm:"not null" json:"filename"`
	OriginalName string `gorm:"not null" json:"original_name"`
	Title        string `json:"title"`
	Description  string `gorm:"type:text" json:"description"`
	Tags         string `json:"tags"`
	Location     string `json:"location"`
	Category     string `json:"category"`

	// EXIF 元数据
	CameraMake   string     `json:"camera_make"`
	CameraModel  string     `json:"camera_model"`
	LensModel    string     `json:"lens_model"`
	FocalLength  string     `json:"focal_length"`
	Aperture     string     `json:"aperture"`
	ShutterSpeed string     `json:"shutter_speed"`
	ISO          string     `json:"iso"`
	DateTaken    *time.Time `json:"date_taken"`

	// 增强 EXIF 元数据（兼容手机和相机拍摄）
	Software     string `json:"software"`      // 拍摄软件/固件版本
	Orientation  int    `json:"orientation"`    // EXIF 方向标记 (1-8)
	WhiteBalance string `json:"white_balance"`  // 白平衡模式
	Flash        string `json:"flash"`          // 闪光灯状态
	ExposureMode string `json:"exposure_mode"`  // 曝光模式
	MeteringMode string `json:"metering_mode"`  // 测光模式
	ExposureBias string `json:"exposure_bias"`  // 曝光补偿
	ColorSpace   string `json:"color_space"`    // 色彩空间
	SceneType    string `json:"scene_type"`     // 场景类型

	// GPS 信息
	GPSLatitude  float64 `json:"gps_latitude"`
	GPSLongitude float64 `json:"gps_longitude"`

	// 文件信息
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`

	// 存储信息
	StorageType   string `json:"storage_type"` // local, s3, minio, aliyun
	StoragePath   string `json:"storage_path"`
	ThumbnailPath string `json:"thumbnail_path"`
	URL           string `gorm:"-" json:"url"`
	ThumbnailURL  string `gorm:"-" json:"thumbnail_url"`

	// Live Photo 支持
	IsLivePhoto   bool   `gorm:"default:false" json:"is_live_photo"`
	LivePhotoPath string `json:"live_photo_path"`
	LivePhotoURL  string `gorm:"-" json:"live_photo_url"`

	// 其他
	Copyright string `json:"copyright"`
	UserID    uint   `json:"user_id"`
	Views     int    `gorm:"default:0" json:"views"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

// PhotoListResponse 照片列表响应
type PhotoListResponse struct {
	Total  int64   `json:"total"`
	Page   int     `json:"page"`
	Size   int     `json:"size"`
	Photos []Photo `json:"photos"`
}

// UploadResponse 上传响应
type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Photo   Photo  `json:"photo"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

