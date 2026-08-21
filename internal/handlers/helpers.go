package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"stayt/internal/cloud"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func generateUniqueKey(userID int, fileName string) string {
	extension := filepath.Ext(fileName)
	key := fmt.Sprintf("users/%d/%s%s", userID, uuid.New().String(), extension)

	return key
}

func generatePresignedURL(ctx context.Context, s3Client *cloud.S3Client, key string) (*string, error) {
	presignClient := s3.NewPresignClient(s3Client.Client)

	presignURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3Client.BucketName),
		Key:    aws.String(key),
	},
		s3.WithPresignExpires(time.Minute*15))

	if err != nil {
		return nil, err
	}

	return &presignURL.URL, nil
}
