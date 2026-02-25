package utils

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

// EXIFData EXIF 数据结构（兼容相机和手机拍摄的照片）
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

	Software     string
	Orientation  int
	WhiteBalance string
	Flash        string
	ExposureMode string
	MeteringMode string
	ExposureBias string
	ColorSpace   string
	SceneType    string
}

// ExtractEXIF 提取 EXIF 信息（兼容相机、iPhone、Android 等设备）
func ExtractEXIF(filePath string) (*EXIFData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	exifData := &EXIFData{
		Width:  img.Width,
		Height: img.Height,
	}

	file.Seek(0, 0)

	x, err := exif.Decode(file)
	if err != nil {
		return exifData, nil
	}

	exifData.CameraMake = getExifString(x, exif.Make)
	exifData.CameraModel = getExifString(x, exif.Model)
	exifData.LensModel = getExifString(x, exif.LensModel)
	exifData.Software = getExifString(x, exif.Software)

	if focalLen, err := x.Get(exif.FocalLength); err == nil {
		if val, err := focalLen.Rat(0); err == nil {
			num, denom := val.Num(), val.Denom()
			if denom.Int64() != 0 {
				fl := float64(num.Int64()) / float64(denom.Int64())
				if fl == float64(int64(fl)) {
					exifData.FocalLength = fmt.Sprintf("%dmm", int64(fl))
				} else {
					exifData.FocalLength = fmt.Sprintf("%.1fmm", fl)
				}
			}
		}
	}

	if fnumber, err := x.Get(exif.FNumber); err == nil {
		if val, err := fnumber.Rat(0); err == nil {
			num, denom := val.Num(), val.Denom()
			if denom.Int64() != 0 {
				f := float64(num.Int64()) / float64(denom.Int64())
				if f == float64(int64(f)) {
					exifData.Aperture = fmt.Sprintf("f/%d", int64(f))
				} else {
					exifData.Aperture = fmt.Sprintf("f/%.1f", f)
				}
			}
		}
	}

	if expTime, err := x.Get(exif.ExposureTime); err == nil {
		if val, err := expTime.Rat(0); err == nil {
			num, denom := val.Num(), val.Denom()
			numInt := num.Int64()
			denomInt := denom.Int64()
			if numInt != 0 && denomInt != 0 {
				if numInt >= denomInt {
					speed := float64(numInt) / float64(denomInt)
					if speed >= 1 {
						exifData.ShutterSpeed = fmt.Sprintf("%.1fs", speed)
					} else {
						exifData.ShutterSpeed = fmt.Sprintf("1/%ds", denomInt/numInt)
					}
				} else {
					exifData.ShutterSpeed = fmt.Sprintf("1/%ds", denomInt/numInt)
				}
			}
		}
	}

	if isoSpeed, err := x.Get(exif.ISOSpeedRatings); err == nil {
		if val, err := isoSpeed.Int(0); err == nil {
			exifData.ISO = fmt.Sprintf("%d", val)
		}
	}

	exifData.DateTaken = parseExifDate(x)

	if lat, long, err := x.LatLong(); err == nil {
		if lat != 0 || long != 0 {
			exifData.GPSLatitude = lat
			exifData.GPSLongitude = long
		}
	}

	if orient, err := x.Get(exif.Orientation); err == nil {
		if val, err := orient.Int(0); err == nil {
			exifData.Orientation = val
			if val >= 5 {
				exifData.Width, exifData.Height = exifData.Height, exifData.Width
			}
		}
	}

	exifData.WhiteBalance = decodeWhiteBalance(x)
	exifData.Flash = decodeFlash(x)
	exifData.ExposureMode = decodeExposureMode(x)
	exifData.MeteringMode = decodeMeteringMode(x)
	exifData.ExposureBias = decodeExposureBias(x)
	exifData.ColorSpace = decodeColorSpace(x)
	exifData.SceneType = decodeSceneType(x)

	return exifData, nil
}

func getExifString(x *exif.Exif, field exif.FieldName) string {
	if tag, err := x.Get(field); err == nil {
		if val, err := tag.StringVal(); err == nil {
			return strings.TrimSpace(strings.TrimRight(val, "\x00"))
		}
	}
	return ""
}

func parseExifDate(x *exif.Exif) *time.Time {
	dateFields := []exif.FieldName{exif.DateTimeOriginal, exif.DateTime}
	formats := []string{
		"2006:01:02 15:04:05",
		"2006-01-02 15:04:05",
		"2006:01:02T15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, field := range dateFields {
		if dateTime, err := x.Get(field); err == nil {
			if val, err := dateTime.StringVal(); err == nil {
				val = strings.TrimSpace(strings.TrimRight(val, "\x00"))
				for _, format := range formats {
					if t, err := time.Parse(format, val); err == nil {
						return &t
					}
				}
			}
		}
	}
	return nil
}

func getExifInt(x *exif.Exif, field exif.FieldName) (int, bool) {
	if tag, err := x.Get(field); err == nil {
		if val, err := tag.Int(0); err == nil {
			return val, true
		}
	}
	return 0, false
}

func decodeWhiteBalance(x *exif.Exif) string {
	if val, ok := getExifInt(x, exif.WhiteBalance); ok {
		switch val {
		case 0:
			return "自动"
		case 1:
			return "手动"
		}
	}
	return ""
}

func decodeFlash(x *exif.Exif) string {
	if val, ok := getExifInt(x, exif.Flash); ok {
		fired := val&1 == 1
		if fired {
			return "已闪光"
		}
		return "未闪光"
	}
	return ""
}

func decodeExposureMode(x *exif.Exif) string {
	if tag, err := x.Get(exif.ExposureProgram); err == nil {
		if val, err := tag.Int(0); err == nil {
			switch val {
			case 1:
				return "手动"
			case 2:
				return "自动"
			case 3:
				return "光圈优先"
			case 4:
				return "快门优先"
			case 5:
				return "创意"
			case 6:
				return "动作"
			case 7:
				return "肖像"
			case 8:
				return "风景"
			}
		}
	}
	return ""
}

func decodeMeteringMode(x *exif.Exif) string {
	if tag, err := x.Get(exif.MeteringMode); err == nil {
		if val, err := tag.Int(0); err == nil {
			switch val {
			case 1:
				return "平均"
			case 2:
				return "中央重点"
			case 3:
				return "点测光"
			case 4:
				return "多点"
			case 5:
				return "评价测光"
			case 6:
				return "局部"
			}
		}
	}
	return ""
}

func decodeExposureBias(x *exif.Exif) string {
	if tag, err := x.Get(exif.ExposureBiasValue); err == nil {
		if val, err := tag.Rat(0); err == nil {
			num, denom := val.Num(), val.Denom()
			if denom.Int64() != 0 {
				bias := float64(num.Int64()) / float64(denom.Int64())
				if bias == 0 {
					return "0 EV"
				}
				return fmt.Sprintf("%+.1f EV", bias)
			}
		}
	}
	return ""
}

func decodeColorSpace(x *exif.Exif) string {
	if val, ok := getExifInt(x, exif.ColorSpace); ok {
		switch val {
		case 1:
			return "sRGB"
		case 2:
			return "Adobe RGB"
		case 0xFFFF:
			return "Uncalibrated"
		}
	}
	return ""
}

func decodeSceneType(x *exif.Exif) string {
	if tag, err := x.Get(exif.SceneCaptureType); err == nil {
		if val, err := tag.Int(0); err == nil {
			switch val {
			case 0:
				return "标准"
			case 1:
				return "风景"
			case 2:
				return "人像"
			case 3:
				return "夜景"
			}
		}
	}
	return ""
}

// GenerateThumbnail 生成缩略图（支持 EXIF 方向自动旋转）
func GenerateThumbnail(srcPath, dstPath string, width, height, quality int) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	img, err := imaging.Open(srcPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("打开图片失败: %w", err)
	}

	thumb := imaging.Fit(img, width, height, imaging.Lanczos)

	if err := imaging.Save(thumb, dstPath, imaging.JPEGQuality(quality)); err != nil {
		return fmt.Errorf("保存缩略图失败: %w", err)
	}

	return nil
}

// GenerateThumbnailBytes 生成缩略图并返回字节数据（用于云存储上传）
func GenerateThumbnailBytes(srcPath string, width, height, quality int) ([]byte, error) {
	img, err := imaging.Open(srcPath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("打开图片失败: %w", err)
	}

	thumb := imaging.Fit(img, width, height, imaging.Lanczos)

	tmpFile, err := os.CreateTemp("", "thumb_*.jpg")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	if err := imaging.Save(thumb, tmpPath, imaging.JPEGQuality(quality)); err != nil {
		return nil, err
	}

	return os.ReadFile(tmpPath)
}

// IsLivePhotoVideo 检查文件是否为 Live Photo 配套视频
func IsLivePhotoVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".mov" || ext == ".mp4"
}

// DetectContentType 根据文件扩展名检测 MIME 类型
func DetectContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	types := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".heic": "image/heic",
		".heif": "image/heif",
		".mov":  "video/quicktime",
		".mp4":  "video/mp4",
	}
	if t, ok := types[ext]; ok {
		return t
	}
	return "application/octet-stream"
}
