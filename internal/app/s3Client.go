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

func (c Client) StreamFile(ctx context.Context,
	bucketName string,
	objectKey string,
	w io.Writer,
) error {

	result, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
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
