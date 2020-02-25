package qqcloud

import (
	"log"
	"net/http"

	"github.com/weihongguo/gglmm"

	cos "github.com/tencentyun/cos-go-sdk-v5"
)

// CosUploadService --
type CosUploadService struct {
	secretID  string
	secretKey string
	region    string
	appID     string
	bucket    string

	cosClient   *cos.Client
	keyFileFunc CosKeyFileFunc
}

// NewCosUploadService --
func NewCosUploadService(secretID string, secretKey string, region string, appID string, bucket string, keyFileFunc CosKeyFileFunc) *CosUploadService {
	cosClient, err := newCosClient(secretID, secretKey, region, appID, bucket)
	if err != nil {
		log.Fatal(err)
	}

	return &CosUploadService{
		secretID:  secretID,
		secretKey: secretKey,
		region:    region,
		appID:     appID,
		bucket:    bucket,

		cosClient:   cosClient,
		keyFileFunc: keyFileFunc,
	}
}

// NewCosUploadServiceConfig --
func NewCosUploadServiceConfig(config ConfigCos, keyFileFunc CosKeyFileFunc) *CosUploadService {
	return NewCosUploadService(config.SecretID, config.SecretKey, config.Region, config.AppID, config.Bucket, keyFileFunc)
}

// CustomActions --
func (service *CosUploadService) CustomActions() ([]*gglmm.HTTPAction, error) {
	actions := []*gglmm.HTTPAction{
		gglmm.NewHTTPAction("/upload", service.Upload, "POST"),
	}
	return actions, nil
}

// RESTAction --
func (service *CosUploadService) RESTAction(action gglmm.RESTAction) (*gglmm.HTTPAction, error) {
	return nil, nil
}

// Upload --
func (service *CosUploadService) Upload(w http.ResponseWriter, r *http.Request) {
	key, file, err := service.keyFileFunc(r)
	if err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	if err = cosPutObj(service.cosClient, key, file); err != nil {
		gglmm.NewFailResponse(err.Error()).WriteJSON(w)
		return
	}

	gglmm.NewSuccessResponse().
		AddData("url", key).
		WriteJSON(w)
}
