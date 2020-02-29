package main

import (
	"mime/multipart"
	"net/http"

	"github.com/weihongguo/gglmm"
	tencentyun "github.com/weihongguo/gglmm-tencentyun"
)

func main() {
	gglmm.RegisterRedisCacher("tcp", "127.0.0.1:6379", 10, 5, 3)
	defer gglmm.CloseRedisCacher()

	gglmm.RegisterBasePath("/api/example")

	// 登录态中间件请参考gglmm-account

	gglmm.RegisterHTTPHandler(tencentyun.NewCosCredentialsService("secretID", "secretKey", "region", "appID", "bucket", cosPrefixKey), "")

	gglmm.RegisterHTTPHandler(tencentyun.NewCosUploadService("secretID", "secretKey", "region", "appID", "bucket", cosKeyFile), "")

	gglmm.ListenAndServe(":10000")
}

func cosPrefixKey(r *http.Request) (string, error) {
	// 根据请求计算路径前缀
	return "example", nil
}

func cosKeyFile(r *http.Request) (string, multipart.File, error) {
	file, _, err := r.FormFile("example")
	if err != nil {
		return "", nil, err
	}

	// 其他判断

	return "example", file, nil
}
