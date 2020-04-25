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

	cosService := tencentyun.NewCosService("secretID", "secretKey", "region", "appID", "bucket").PrefixKeyFunc(cosPrefixKey).KeyFileFunc(cosKeyFile)

	gglmm.HandleHTTPAction("/cos/credentials", cosService.Credentials, "GET")

	gglmm.HandleHTTPAction("/cos/upload", cosService.Upload, "POST")

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
