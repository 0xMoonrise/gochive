package core

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	cfg "github.com/0xMoonrise/gochive/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type clientS3 struct {
	S3Client *s3.Client
}

func (c *clientS3) GetItem(ctx context.Context, objKey string) (obj *Object, err error) {
	result, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(cfg.BUCKET),
		Key:    aws.String(objKey),
	})
	if err != nil {
		return nil, err
	}

	obj = &Object{Reader: result.Body}
	obj.Length = int64(0)
	if result.ContentLength != nil {
		obj.Length = *result.ContentLength
	}
	obj.ContentType = "application/octet-stream"
	if result.ContentType != nil {
		obj.ContentType = *result.ContentType
	}
	return obj, nil
}

func (c *clientS3) PutItem(ctx context.Context, objKey string, obj *Object) (err error) {
	data, err := io.ReadAll(obj.Reader)
	if err != nil {
		return err
	}
	_, err = c.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(cfg.BUCKET),
		Key:           aws.String(objKey),
		Body:          bytes.NewReader(data),
		ContentLength: &obj.Length,
		ContentType:   aws.String(obj.ContentType),
	})
	return
}

func (c *clientS3) DelItem(ctx context.Context, objKey string) (err error) {
	_, err = c.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.BUCKET),
		Key:    aws.String(objKey),
	})
	return
}

func NewS3Client() (*clientS3, error) {

	if cfg.ACCESS_KEY == "" || cfg.SECRET_KEY == "" {
		return nil, errors.New("no credentials were provided")
	}
	creds := credentials.NewStaticCredentialsProvider(cfg.ACCESS_KEY, cfg.SECRET_KEY, "")
	conf, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(cfg.REGION),
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

	client := s3.NewFromConfig(conf, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3_ENDPOINT)
		o.UsePathStyle = true
		o.EndpointOptions.DisableHTTPS = true
	})

	return &clientS3{
		S3Client: client,
	}, nil
}
