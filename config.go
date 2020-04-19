package tencentyun

import (
	"strings"
)

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

// ConfigTencentYun --
type ConfigTencentYun struct {
	Cos ConfigCos
}

// Check --
func (config ConfigTencentYun) Check() bool {
	if !config.Cos.Check() {
		return false
	}
	return true
}
