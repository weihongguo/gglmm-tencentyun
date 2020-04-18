package main

import (
	"mime/multipart"
	"net/http"

	"github.com/weihongguo/gglmm"
	redis "github.com/weihongguo/gglmm-redis"
	tencentyun "github.com/weihongguo/gglmm-tencentyun"
)

func main() {
	redisCacher := redis.NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 30)
	defer redisCacher.Close()
	gglmm.RegisterCacher(redisCacher)

	gglmm.BasePath("/api/example")

	// 登录态中间件请参考gglmm-account

	gglmm.HandleHTTP(tencentyun.NewCosCredentialsService("secretID", "secretKey", "region", "appID", "bucket", cosPrefixKey), "")

	gglmm.HandleHTTP(tencentyun.NewCosUploadService("secretID", "secretKey", "region", "appID", "bucket", cosKeyFile), "")

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
