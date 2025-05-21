package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Store struct {
	s3 *minio.Client
}

func New(creds Credentials) (*Store, error) {
	minioClient, err := minio.New(creds.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(creds.accessKeyID, creds.secretAccessKey, ""),
		Secure: creds.useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 storage: %w", err)
	}
	return &Store{
		s3: minioClient,
	}, nil
}

func (s Store) GetBaseURL() string {
	return s.s3.EndpointURL().String()
}

func (s Store) CreateBuckets(ctx context.Context, bucketName ...string) error {
	for _, name := range bucketName {
		err := s.s3.MakeBucket(ctx, name, minio.MakeBucketOptions{})
		if err != nil {
			if exists, errBucketExists := s.s3.BucketExists(ctx, name); errBucketExists == nil && exists {
				continue
			}
			return fmt.Errorf("failed to create bucket %s: %w", name, err)
		}
	}
	return nil
}

// Upload uploads an object to S3.
func (s Store) Upload(ctx context.Context, object *Object) error {
	options := minio.PutObjectOptions{}
	if object.metadata != nil {
		options.UserMetadata = object.metadata
	}
	if object.contentType != "" {
		options.ContentType = string(object.contentType)
	}

	_, err := s.s3.PutObject(ctx,
		object.bucketName,
		object.Key,
		bytes.NewReader(object.data),
		int64(len(object.data)),
		options,
	)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}
	return nil
}

func (s Store) Download(ctx context.Context, bucketName, key string) (*Object, error) {
	obj, err := s.s3.GetObject(ctx, bucketName, key, minio.GetObjectOptions{
		ServerSideEncryption: nil,
		PartNumber:           0,
		Checksum:             false,
		Internal:             minio.AdvancedGetOptions{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	objInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat object: %w", err)
	}

	defer obj.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(obj); err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}
	return &Object{
		bucketName:  bucketName,
		Key:         key,
		data:        buf.Bytes(),
		contentType: ContentType(objInfo.ContentType),
		metadata:    objInfo.UserMetadata,
	}, nil
}

// TODO: add listObjects method
// TODO: add upload multiple method
