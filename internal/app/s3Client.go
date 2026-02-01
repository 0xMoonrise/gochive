package app

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	S3Client *s3.Client
}

func (c *Client) GetItem(ctx context.Context, objKey string) (
	length int64,
	contentType string,
	reader io.ReadCloser,
	err error,
) {
	bucket := os.Getenv("BUCKET")

	result, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objKey),
	})

	if err != nil {
		return
	}

	length = int64(0)
	if result.ContentLength != nil {
		length = *result.ContentLength
	}

	contentType = "application/octet-stream"
	if result.ContentType != nil {
		contentType = *result.ContentType
	}

	reader = result.Body
	return
}

func NewS3Client() (*Client, error) {

	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	endpoint := os.Getenv("S3_ENDPOINT")
	region := os.Getenv("REGION")

	if endpoint == "" {
		endpoint = "http://localhost:3901"
	}

	if accessKey == "" || secretKey == "" {
		return nil, errors.New("No credentials were provided")
	}

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &Client{
		S3Client: client,
	}, nil
}
