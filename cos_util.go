package tencentyun

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	cos "github.com/tencentyun/cos-go-sdk-v5"
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

func newStsClient(secretID string, secretKey string) (*sts.Client, error) {
	return sts.NewClient(secretID, secretKey, nil), nil
}

// CosPrefixKeyFunc --
type CosPrefixKeyFunc func(r *http.Request) (config *ConfigCos, prefixKey string, err error)

func stsGetCredential(stsClient *sts.Client, region string, appID string, bucket string, prefixKey string) (*sts.CredentialResult, error) {
	opts := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:" + region + ":uid/" + appID + ":" + bucket + "/" + prefixKey + "/*",
					},
				},
			},
		},
	}

	return stsClient.GetCredential(opts)
}

func newCosClient(secretID string, secretKey string, region string, appID string, bucket string) (*cos.Client, error) {
	bucketURL, err := url.Parse("https://" + bucket + ".cos." + region + ".myqcloud.com")
	if err != nil {
		return nil, err
	}
	baseURL := &cos.BaseURL{
		BucketURL: bucketURL,
	}
	cosClient := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return cosClient, nil
}

// CosKeyFileFunc --
type CosKeyFileFunc func(r *http.Request) (config *ConfigCos, key string, file multipart.File, err error)

func cosPutObj(cosClient *cos.Client, key string, file multipart.File) error {
	res, err := cosClient.Object.Put(context.Background(), key, file, nil)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return ErrCosUpload
	}
	return nil
}
