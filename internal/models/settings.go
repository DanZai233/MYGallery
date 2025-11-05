package models

import (
	"time"

	"gorm.io/gorm"
)

// Settings 网站设置模型
type Settings struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 网站基本信息
	SiteTitle       string `json:"site_title"`
	SiteDescription string `gorm:"type:text" json:"site_description"`
	SiteKeywords    string `json:"site_keywords"`
	SiteLogo        string `json:"site_logo"`
	
	// 页眉页脚
	HeaderHTML  string `gorm:"type:text" json:"header_html"`
	FooterHTML  string `gorm:"type:text" json:"footer_html"`
	
	// 备案信息
	ICPRecord     string `json:"icp_record"`      // ICP备案号
	BeianRecord   string `json:"beian_record"`    // 公安备案号
	BeianLink     string `json:"beian_link"`      // 公安备案链接
	
	// 联系方式
	ContactEmail  string `json:"contact_email"`
	ContactPhone  string `json:"contact_phone"`
	ContactWeChat string `json:"contact_wechat"`
	
	// 社交媒体
	GithubURL    string `json:"github_url"`
	TwitterURL   string `json:"twitter_url"`
	WeiboURL     string `json:"weibo_url"`
	
	// 统计代码
	AnalyticsCode string `gorm:"type:text" json:"analytics_code"`
	
	// 其他设置
	CustomCSS     string `gorm:"type:text" json:"custom_css"`
	CustomJS      string `gorm:"type:text" json:"custom_js"`
}

// Category 分类模型
type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Cover       string `json:"cover"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	PhotoCount  int    `gorm:"default:0" json:"photo_count"`
}

