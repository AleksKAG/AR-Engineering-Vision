package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	c      *minio.Client
	bucket string
}

func NewClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*Client, error) {
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	// ensure bucket exists
	ctx := context.Background()
	exists, err := mc.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := mc.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make bucket: %w", err)
		}
	}
	return &Client{c: mc, bucket: bucket}, nil
}

func (s *Client) Upload(ctx context.Context, objectName string, r io.Reader, size int64, contentType string) (string, error) {
	_, err := s.c.PutObject(ctx, s.bucket, objectName, r, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}
	// return object location (for now objectName)
	return objectName, nil
}

func (s *Client) URL(objectName string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presigned, err := s.c.PresignedGetObject(context.Background(), s.bucket, objectName, expires, reqParams)
	if err != nil {
		return "", err
	}
	return presigned.String(), nil
}
