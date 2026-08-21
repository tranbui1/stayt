package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client     *s3.Client
	BucketName string
}

// TODO: route client to pass into handlers/thoughts.go file, call os to get bucketName
func Config(bucketName string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	cloudCFG := &S3Client{
		Client:     s3Client,
		BucketName: bucketName,
	}
	return cloudCFG, nil
}
