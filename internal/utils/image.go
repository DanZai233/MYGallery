package utils

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// EXIFData EXIF 数据结构
type EXIFData struct {
	CameraMake   string
	CameraModel  string
	LensModel    string
	FocalLength  string
	Aperture     string
	ShutterSpeed string
	ISO          string
	DateTaken    *time.Time
	GPSLatitude  float64
	GPSLongitude float64
	Width        int
	Height       int
}

// ExtractEXIF 提取 EXIF 信息
func ExtractEXIF(filePath string) (*EXIFData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 解析图片获取尺寸
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}
	
	exifData := &EXIFData{
		Width:  img.Width,
		Height: img.Height,
	}
	
	// 重置文件指针
	file.Seek(0, 0)
	
	// 解析 EXIF
	x, err := exif.Decode(file)
	if err != nil {
		// 如果没有 EXIF 数据，返回基本信息
		return exifData, nil
	}
	
	// 相机制造商
	if make, err := x.Get(exif.Make); err == nil {
		if val, err := make.StringVal(); err == nil {
			exifData.CameraMake = val
		}
	}
	
	// 相机型号
	if model, err := x.Get(exif.Model); err == nil {
		if val, err := model.StringVal(); err == nil {
			exifData.CameraModel = val
		}
	}
	
	// 镜头型号
	if lens, err := x.Get(exif.LensModel); err == nil {
		if val, err := lens.StringVal(); err == nil {
			exifData.LensModel = val
		}
	}
	
	// 焦距
	if focalLen, err := x.Get(exif.FocalLength); err == nil {
		if val, err := focalLen.Rat(0); err == nil {
			num, denom := val.Num(), val.Denom()
			if denom.Int64() != 0 {
				exifData.FocalLength = fmt.Sprintf("%.0fmm", float64(num.Int64())/float64(denom.Int64()))
			}
		}
	}
	
	// 光圈
	if fnumber, err := x.Get(exif.FNumber); err == nil {
		if val, err := fnumber.Rat(0); err == nil {
			num, denom := val.Num(), val.Denom()
			if denom.Int64() != 0 {
				exifData.Aperture = fmt.Sprintf("f/%.1f", float64(num.Int64())/float64(denom.Int64()))
			}
		}
	}
	
	// 快门速度
	if expTime, err := x.Get(exif.ExposureTime); err == nil {
		if val, err := expTime.Rat(0); err == nil {
			num, denom := val.Num(), val.Denom()
			numInt := num.Int64()
			denomInt := denom.Int64()
			if numInt != 0 && denomInt != 0 {
				if numInt >= denomInt {
					exifData.ShutterSpeed = fmt.Sprintf("%.1fs", float64(numInt)/float64(denomInt))
				} else {
					exifData.ShutterSpeed = fmt.Sprintf("1/%ds", denomInt/numInt)
				}
			}
		}
	}
	
	// ISO
	if isoSpeed, err := x.Get(exif.ISOSpeedRatings); err == nil {
		if val, err := isoSpeed.Int(0); err == nil {
			exifData.ISO = fmt.Sprintf("%d", val)
		}
	}
	
	// 拍摄时间
	if dateTime, err := x.Get(exif.DateTimeOriginal); err == nil {
		if val, err := dateTime.StringVal(); err == nil {
			if t, err := time.Parse("2006:01:02 15:04:05", val); err == nil {
				exifData.DateTaken = &t
			}
		}
	}
	
	// GPS 信息
	if lat, long, err := x.LatLong(); err == nil {
		exifData.GPSLatitude = lat
		exifData.GPSLongitude = long
	}
	
	return exifData, nil
}

// GenerateThumbnail 生成缩略图
func GenerateThumbnail(srcPath, dstPath string, width, height, quality int) error {
	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}
	
	// 打开源图片
	img, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开图片失败: %w", err)
	}
	
	// 生成缩略图
	thumb := imaging.Fit(img, width, height, imaging.Lanczos)
	
	// 保存缩略图
	if err := imaging.Save(thumb, dstPath, imaging.JPEGQuality(quality)); err != nil {
		return fmt.Errorf("保存缩略图失败: %w", err)
	}
	
	return nil
}

// GetRational 辅助函数：从 tiff.Tag 获取有理数（已弃用，仅保留作为参考）
// 注意：rat.Num() 和 rat.Denom() 返回的是 *big.Int 类型，需要使用 .Int64() 方法转换

