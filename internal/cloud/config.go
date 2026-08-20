package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client     *s3.Client
	TmClient   *transfermanager.Client
	BucketName string
}

// TODO: route client to pass into handlers/thoughts.go file, call os to get bucketName
func Config(bucketName string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	tmClient := transfermanager.New(s3Client, func(o *transfermanager.Options) {
		o.Concurrency = 5
		o.PartSizeBytes = 5 * 1024 * 1024
	})

	cloudCFG := &S3Client{
		Client:     s3Client,
		TmClient:   tmClient,
		BucketName: bucketName,
	}
	return cloudCFG, nil
}
