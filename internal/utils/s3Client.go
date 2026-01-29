package utils

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	S3Client *s3.Client
}

func (b Bucket) StreamFile(
	ctx context.Context,
	bucketName string,
	objectKey string,
	w io.Writer,
) error {

	result, err := b.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return err
	}

	defer result.Body.Close()

	_, err = io.Copy(w, result.Body)
	return err
}

func NewS3Client() *Bucket {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal(err)
	}
	return &Bucket{
		S3Client: s3.NewFromConfig(sdkConfig),
	}
}
