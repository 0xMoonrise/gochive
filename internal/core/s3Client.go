package core

import (
	"context"
	"errors"
	"net/http"
	"time"

	cfg "github.com/0xMoonrise/gochive/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type clientS3 struct {
	S3Client *s3.Client
	manager  *transfermanager.Client
	conf     *cfg.S3ClientConfig
}

func (c *clientS3) GetItem(ctx context.Context, objKey string) (obj *Object, err error) {
	result, err := c.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.conf.Bucket),
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

	return
}

func (c *clientS3) PutItem(ctx context.Context, objKey string, obj *Object) (err error) {
	_, err = c.manager.UploadObject(ctx, &transfermanager.UploadObjectInput{
		Bucket:      aws.String(c.conf.Bucket),
		Key:         aws.String(objKey),
		Body:        obj.Reader,
		ContentType: aws.String(obj.ContentType),
	})
	return
}

func (c *clientS3) DelItem(ctx context.Context, objKey string) (err error) {
	_, err = c.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.conf.Bucket),
		Key:    aws.String(objKey),
	})
	return
}

func (app App) NewS3Client() (*clientS3, error) {

	if app.Config.S3.AccessKey == "" || app.Config.S3.SecretKey == "" {
		return nil, errors.New("no credentials were provided")
	}

	creds := credentials.NewStaticCredentialsProvider(
		app.Config.S3.AccessKey,
		app.Config.S3.SecretKey,
		"",
	)

	conf, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(app.Config.S3.Region),
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
		o.BaseEndpoint = aws.String(app.Config.S3.S3Endpoint)
		o.UsePathStyle = true
		o.EndpointOptions.DisableHTTPS = true
		o.DisableLogOutputChecksumValidationSkipped = true
	})

	_, err = client.HeadBucket(
		context.Background(),
		&s3.HeadBucketInput{
			Bucket: &app.Config.S3.Bucket,
		})

	if err != nil {
		return nil, err
	}

	manager := transfermanager.New(client, func(o *transfermanager.Options) {
		o.PartSizeBytes = 64 * 1024 * 1024
		o.Concurrency = 3
	})

	return &clientS3{
		S3Client: client,
		manager:  manager,
		conf:     &app.Config.S3,
	}, nil
}
