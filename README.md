# gglmm-tencentyun
## 依赖
+ github.com/tencentyun/cos-go-sdk-v5
+ github.com/tencentyun/qcloud-cos-sts-sdk
## COS
```golang
func NewCosServiceConfig(config ConfigCos) *CosService
func NewCosService(secretID string, secretKey string, region string, appID string, bucket string) *CosService
```
## 使用方法
+ 参考example