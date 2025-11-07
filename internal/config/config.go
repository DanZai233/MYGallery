package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Storage  StorageConfig  `yaml:"storage"`
	JWT      JWTConfig      `yaml:"jwt"`
	Admin    AdminConfig    `yaml:"admin"`
	Image    ImageConfig    `yaml:"image"`
	CORS     CORSConfig     `yaml:"cors"`
	App      AppConfig      `yaml:"app"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Type     string         `yaml:"type"`
	SQLite   SQLiteConfig   `yaml:"sqlite"`
	MySQL    MySQLConfig    `yaml:"mysql"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type StorageConfig struct {
	Type   string        `yaml:"type"`
	Local  LocalStorage  `yaml:"local"`
	S3     S3Storage     `yaml:"s3"`
	MinIO  MinIOStorage  `yaml:"minio"`
	Aliyun AliyunStorage `yaml:"aliyun"`
}

type LocalStorage struct {
	UploadDir    string `yaml:"upload_dir"`
	ThumbnailDir string `yaml:"thumbnail_dir"`
	URLPrefix    string `yaml:"url_prefix"`
}

type S3Storage struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
	URLPrefix string `yaml:"url_prefix"`
}

type MinIOStorage struct {
	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	UseSSL    bool   `yaml:"use_ssl"`
	URLPrefix string `yaml:"url_prefix"`
}

type AliyunStorage struct {
	Endpoint  string `yaml:"endpoint"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	URLPrefix string `yaml:"url_prefix"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}

type ImageConfig struct {
	MaxUploadSize int64    `yaml:"max_upload_size"`
	AllowedTypes  []string `yaml:"allowed_types"`
	Thumbnail     struct {
		Width   int `yaml:"width"`
		Height  int `yaml:"height"`
		Quality int `yaml:"quality"`
	} `yaml:"thumbnail"`
}

type CORSConfig struct {
	Enabled      bool     `yaml:"enabled"`
	AllowOrigins []string `yaml:"allow_origins"`
	AllowMethods []string `yaml:"allow_methods"`
	AllowHeaders []string `yaml:"allow_headers"`
}

type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	SiteTitle   string `yaml:"site_title"`
	Description string `yaml:"site_description"`
	Pagination  struct {
		PageSize int `yaml:"page_size"`
	} `yaml:"pagination"`
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return FromBytes(data)
}

// FromBytes parses configuration data from YAML bytes.
func FromBytes(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			Mode: "release",
		},
		Database: DatabaseConfig{
			Type: "sqlite",
			SQLite: SQLiteConfig{
				Path: "./data/mygallery.db",
			},
		},
		Storage: StorageConfig{
			Type: "local",
			Local: LocalStorage{
				UploadDir:    "./uploads",
				ThumbnailDir: "./uploads/thumbnails",
				URLPrefix:    "/uploads",
			},
		},
		JWT: JWTConfig{
			Secret:      "change-this-secret-key",
			ExpireHours: 168,
		},
		Admin: AdminConfig{
			Username: "admin",
			Password: "admin123",
			Email:    "admin@example.com",
		},
		Image: ImageConfig{
			MaxUploadSize: 52428800,
			AllowedTypes:  []string{"image/jpeg", "image/png", "image/gif", "image/webp"},
		},
		CORS: CORSConfig{
			Enabled:      true,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		},
		App: AppConfig{
			Name:        "MYGallery",
			Version:     "2.0.0",
			SiteTitle:   "我的照片墙",
			Description: "个人照片展示与管理系统",
		},
	}
}
