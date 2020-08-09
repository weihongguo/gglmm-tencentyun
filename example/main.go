package main

import (
	"mime/multipart"
	"net/http"

	"github.com/weihongguo/gglmm"
	redis "github.com/weihongguo/gglmm-redis"
	tencentyun "github.com/weihongguo/gglmm-tencentyun"
)

func main() {
	cacher := redis.NewCacher("tcp", "127.0.0.1:6379", 5, 10, 3, 30)
	defer cacher.Close()

	gglmm.BasePath("/api/example")

	cosService := tencentyun.NewCosService(cosPrefixKey, cosKeyFile, cacher)

	gglmm.HandleHTTPAction("/cos/credentials", cosService.Credentials, "GET")

	gglmm.HandleHTTPAction("/cos/upload", cosService.Upload, "POST")

	gglmm.ListenAndServe(":10000")
}

func cosPrefixKey(r *http.Request) (*tencentyun.ConfigCos, string, error) {
	// TODO
	// 根据请求计算路径前缀
	return nil, "example", nil
}

func cosKeyFile(r *http.Request) (*tencentyun.ConfigCos, string, multipart.File, error) {
	file, _, err := r.FormFile("example")
	if err != nil {
		return nil, "", nil, err
	}
	// TODO
	// 其他判断
	return nil, "example", file, nil
}
