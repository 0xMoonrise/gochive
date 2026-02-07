package core

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	S3Client *s3.Client
}

func (c *Client) GetItem(ctx context.Context, objKey string) (
	obj *Object,
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

	obj = &Object{}
	obj.Length = int64(0)
	if result.ContentLength != nil {
		obj.Length = *result.ContentLength
	}

	obj.ContentType = "application/octet-stream"
	if result.ContentType != nil {
		obj.ContentType = *result.ContentType
	}

	obj.Reader = result.Body
	return
}

func (c *Client) PutItem(
	ctx context.Context,
	objKey string,
	obj *Object,
) (
	err error,
) {
	bucket := os.Getenv("BUCKET")
	_, err = c.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objKey),
		Body:          obj.Reader,
		ContentLength: &obj.Length,
		ContentType:   aws.String(obj.ContentType),
	})
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
		config.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       300 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableCompression:    true,
			},
		}))

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		o.EndpointOptions.DisableHTTPS = true
	})

	return &Client{
		S3Client: client,
	}, nil
}
