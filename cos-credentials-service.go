package qqcloud

import (
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

// CosCredentialsService --
type CosCredentialsService struct {
	secretID      string
	secretKey     string
	region        string
	appID         string
	bucket        string
	stsClient     *sts.Client
	prefixKeyFunc CosPrefixKeyFunc
}

// NewCosCredentialsService --
func NewCosCredentialsService(secretID string, secretKey string, region string, appID string, bucket string, prefixKeyFunc CosPrefixKeyFunc) *CosCredentialsService {
	stsClient, err := newStsClient(secretID, secretKey)
	if err != nil {
		log.Fatal(err)
	}

	return &CosCredentialsService{
		secretID:      secretID,
		secretKey:     secretKey,
		region:        region,
		appID:         appID,
		bucket:        bucket,
		stsClient:     stsClient,
		prefixKeyFunc: prefixKeyFunc,
	}
}

// NewCosCredentialsServiceConfig --
func NewCosCredentialsServiceConfig(config ConfigCos, prefixKeyFunc CosPrefixKeyFunc) *CosCredentialsService {
	return NewCosCredentialsService(config.SecretID, config.SecretKey, config.Region, config.AppID, config.Bucket, prefixKeyFunc)
}

// CustomActions --
func (service *CosCredentialsService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/credentials", service.Credentials, "GET"),
	}
	return actions, nil
}

// RESTAction --
func (service *CosCredentialsService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// Credentials --
func (service *CosCredentialsService) Credentials(w http.ResponseWriter, r *http.Request) {
	prefixKey, err := service.prefixKeyFunc(r)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	res := &sts.CredentialResult{}
	cacher := gglmm.DefaultCacher()
	if cacher != nil {
		if err := cacher.GetObj("cos-credentials-"+prefixKey, res); err == nil {
			gglmm.NewSuccessResponse().
				AddData("credentials", res.Credentials).
				AddData("expiredTime", res.ExpiredTime).
				AddData("expiration", res.Expiration).
				WriteJSON(w)
			return
		}
	}

	res, err = stsGetCredential(service.stsClient, service.region, service.appID, service.bucket, prefixKey)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	if cacher != nil {
		cacher.SetEx("cos-credentials-"+prefixKey, res, 30*60)
	}

	gglmm.NewSuccessResponse().
		AddData("credentials", res.Credentials).
		AddData("expiredTime", res.ExpiredTime).
		AddData("expiration", res.Expiration).
		WriteJSON(w)
}
