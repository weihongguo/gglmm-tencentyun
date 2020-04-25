package tencentyun

import (
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// CosService --
type CosService struct {
	secretID  string
	secretKey string
	region    string
	appID     string
	bucket    string

	stsClient     *sts.Client
	prefixKeyFunc CosPrefixKeyFunc

	cosClient   *cos.Client
	keyFileFunc CosKeyFileFunc
}

// NewCosService --
func NewCosService(secretID string, secretKey string, region string, appID string, bucket string) *CosService {
	return &CosService{
		secretID:  secretID,
		secretKey: secretKey,
		region:    region,
		appID:     appID,
		bucket:    bucket,
	}
}

// NewCosServiceConfig --
func NewCosServiceConfig(config ConfigCos) *CosService {
	return NewCosService(config.SecretID, config.SecretKey, config.Region, config.AppID, config.Bucket)
}

// PrefixKeyFunc --
func (service *CosService) PrefixKeyFunc(prefixKeyFunc CosPrefixKeyFunc) *CosService {
	stsClient, err := newStsClient(service.secretID, service.secretKey)
	if err != nil {
		log.Fatal(err)
	}
	service.stsClient = stsClient
	service.prefixKeyFunc = prefixKeyFunc
	return service
}

// KeyFileFunc --
func (service *CosService) KeyFileFunc(keyFileFunc CosKeyFileFunc) *CosService {
	cosClient, err := newCosClient(service.secretID, service.secretKey, service.region, service.appID, service.bucket)
	if err != nil {
		log.Fatal(err)
	}
	service.cosClient = cosClient
	service.keyFileFunc = keyFileFunc
	return service
}

// Credentials --
func (service *CosService) Credentials(w http.ResponseWriter, r *http.Request) {
	prefixKey, err := service.prefixKeyFunc(r)
	if err != nil {
		gglmm.ErrorResponse(err.Error()).JSON(w)
		return
	}

	result := &sts.CredentialResult{}
	cacher := gglmm.DefaultCacher()
	if cacher != nil {
		if err := cacher.GetObj("cos-credentials-"+prefixKey, result); err == nil {
			gglmm.OkResponse().
				AddData("credentials", result.Credentials).
				AddData("expiredTime", result.ExpiredTime).
				AddData("expiration", result.Expiration).
				JSON(w)
			return
		}
	}

	result, err = stsGetCredential(service.stsClient, service.region, service.appID, service.bucket, prefixKey)
	if err != nil {
		gglmm.ErrorResponse(err.Error()).JSON(w)
		return
	}

	if cacher != nil {
		cacher.SetEx("cos-credentials-"+prefixKey, result, 30*60)
	}

	gglmm.OkResponse().
		AddData("credentials", result.Credentials).
		AddData("expiredTime", result.ExpiredTime).
		AddData("expiration", result.Expiration).
		JSON(w)
}

// Upload --
func (service *CosService) Upload(w http.ResponseWriter, r *http.Request) {
	key, file, err := service.keyFileFunc(r)
	if err != nil {
		gglmm.ErrorResponse(err.Error()).JSON(w)
		return
	}

	if err = cosPutObj(service.cosClient, key, file); err != nil {
		gglmm.ErrorResponse(err.Error()).JSON(w)
		return
	}

	gglmm.OkResponse().
		AddData("url", key).
		JSON(w)
}
