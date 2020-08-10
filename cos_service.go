package tencentyun

import (
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// CosService --
type CosService struct {
	prefixKeyFunc CosPrefixKeyFunc
	keyFileFunc   CosKeyFileFunc
	cacher        gglmm.Cacher
}

// NewCosService --
func NewCosService(prefixKeyFunc CosPrefixKeyFunc, keyFileFunc CosKeyFileFunc, cacher gglmm.Cacher) *CosService {
	return &CosService{
		prefixKeyFunc: prefixKeyFunc,
		keyFileFunc:   keyFileFunc,
		cacher:        cacher,
	}
}

// Credentials --
func (service *CosService) Credentials(w http.ResponseWriter, r *http.Request) {
	config, prefixKey, err := service.prefixKeyFunc(r)
	if err != nil {
		gglmm.FailResponse(gglmm.NewErrFileLine(err)).JSON(w)
		return
	}
	stsClient, err := newStsClient(config.SecretID, config.SecretKey)
	if err != nil {
		log.Fatal(err)
	}

	result := &sts.CredentialResult{}
	if service.cacher != nil {
		if err := service.cacher.GetObj("cos-credentials-"+prefixKey, result); err == nil {
			gglmm.OkResponse().
				AddData("credentials", result.Credentials).
				AddData("expiredTime", result.ExpiredTime).
				AddData("expiration", result.Expiration).
				JSON(w)
			return
		}
	}

	result, err = stsGetCredential(stsClient, config.Region, config.AppID, config.Bucket, prefixKey)
	if err != nil {
		gglmm.FailResponse(gglmm.NewErrFileLine(err)).JSON(w)
		return
	}

	if service.cacher != nil {
		service.cacher.SetEx("cos-credentials-"+prefixKey, result, 30*60)
	}

	gglmm.OkResponse().
		AddData("credentials", result.Credentials).
		AddData("expiredTime", result.ExpiredTime).
		AddData("expiration", result.Expiration).
		JSON(w)
}

// Upload --
func (service *CosService) Upload(w http.ResponseWriter, r *http.Request) {
	config, key, file, err := service.keyFileFunc(r)
	if err != nil {
		gglmm.FailResponse(gglmm.NewErrFileLine(err)).JSON(w)
		return
	}
	cosClient, err := newCosClient(config.SecretID, config.SecretKey, config.Region, config.AppID, config.Bucket)
	if err != nil {
		log.Fatal(err)
	}
	if err = cosPutObj(cosClient, key, file); err != nil {
		gglmm.FailResponse(gglmm.NewErrFileLine(err)).JSON(w)
		return
	}
	gglmm.OkResponse().
		AddData("url", key).
		JSON(w)
}
