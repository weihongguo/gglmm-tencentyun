package tencentyun

import (
	"errors"
	"strings"
)

// ErrCosUpload --
var ErrCosUpload = errors.New("Cos上传文件失败")

// ConfigCos --
type ConfigCos struct {
	SecretID  string
	SecretKey string
	Region    string
	AppID     string
	Bucket    string
}

// Check --
func (config ConfigCos) Check() bool {
	if config.SecretID == "" || config.SecretKey == "" {
		return false
	}
	if config.Region == "" {
		return false
	}
	if config.AppID == "" || config.Bucket == "" || !strings.Contains(config.Bucket, "-") {
		return false
	}
	return true
}

// ConfigTencentyun --
type ConfigTencentyun struct {
	Cos ConfigCos
}
