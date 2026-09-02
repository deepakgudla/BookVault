package providers

import (
	"context"
	"log"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appConfig "github.com/deepakgudla/bookvault/internal/config"
)

// S3Provider stores uploaded files in an S3-compatible object store.
type S3Provider struct {
	client   *s3.Client
	uploader *manager.Uploader //nolint:staticcheck // transfermanager is not compatible with this SDK version.
	bucket   string
	endpoint string
}

// NewS3Provider creates an S3 upload provider.
func NewS3Provider(cfg *appConfig.Config) *S3Provider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKey,
			cfg.AWS.SecretKey,
			"",
		)),
	)

	if err != nil {
		panic("failed to create AWS config" + err.Error())
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.AWS.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.AWS.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Provider{
		client:   client,
		uploader: manager.NewUploader(client), //nolint:staticcheck // transfermanager is not compatible with this SDK version.
		bucket:   cfg.AWS.S3Bucket,
		endpoint: cfg.AWS.S3Endpoint,
	}
}

// UploadFile uploads a multipart file to object storage.
func (p *S3Provider) UploadFile(file *multipart.FileHeader, path string) (string, error) {

	log.Printf("uploading file %s using S3", path)
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("failed to close source file: %v", err)
		}
	}()

	result, err := p.uploader.Upload(context.TODO(), &s3.PutObjectInput{ //nolint:staticcheck // transfermanager is not compatible with this SDK version.
		Bucket: aws.String(p.bucket),
		Key:    aws.String(path),
		Body:   src,
	})

	if err != nil {
		return "", err
	}

	return *result.Key, nil
}

// DeleteFile removes a file from object storage.
func (p *S3Provider) DeleteFile(path string) error {
	_, err := p.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	})

	return err
}
