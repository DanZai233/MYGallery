package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mygallery/mygallery/internal/config"
	"github.com/mygallery/mygallery/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) error {
	var err error
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "sqlite":
		// 确保数据目录存在
		dir := filepath.Dir(cfg.Database.SQLite.Path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建数据目录失败: %w", err)
		}
		dialector = sqlite.Open(cfg.Database.SQLite.Path)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.Database.MySQL.Username,
			cfg.Database.MySQL.Password,
			cfg.Database.MySQL.Host,
			cfg.Database.MySQL.Port,
			cfg.Database.MySQL.Database,
			cfg.Database.MySQL.Charset,
		)
		dialector = mysql.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Postgres.Host,
			cfg.Database.Postgres.Port,
			cfg.Database.Postgres.Username,
			cfg.Database.Postgres.Password,
			cfg.Database.Postgres.Database,
			cfg.Database.Postgres.SSLMode,
		)
		dialector = postgres.Open(dsn)

	default:
		return fmt.Errorf("不支持的数据库类型: %s", cfg.Database.Type)
	}

	// 设置日志级别
	logLevel := logger.Silent
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	// 自动迁移表结构
	if err := DB.AutoMigrate(&models.User{}, &models.Photo{}, &models.Settings{}, &models.Category{}, &models.PhotoReaction{}, &models.Album{}, &models.PhotoAlbum{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// CreateDefaultAdmin 创建默认管理员账号
func CreateDefaultAdmin(cfg *config.Config) error {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	
	if count > 0 {
		// 已有用户，不创建
		return nil
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	admin := models.User{
		Username: cfg.Admin.Username,
		Password: string(hashedPassword),
		Email:    cfg.Admin.Email,
		Role:     "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("创建管理员失败: %w", err)
	}

	log.Println("✓ 默认管理员账号已创建")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

