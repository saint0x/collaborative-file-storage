package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type B2Service struct {
	client     *s3.Client
	bucketName string
	region     string
}

func NewB2Service(keyID, applicationKey, bucketName, endpoint, region string) (*B2Service, error) {
	b2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s", endpoint),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(b2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(keyID, applicationKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load Backblaze B2 configuration: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &B2Service{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

func (b *B2Service) UploadFile(ctx context.Context, key string, content io.Reader) error {
	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
		Body:   content,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	return nil
}

func (b *B2Service) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	return result.Body, nil
}

func (b *B2Service) DeleteFile(ctx context.Context, key string) error {
	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

func (b *B2Service) GetSignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(b.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	return request.URL, nil
}

func (b *B2Service) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	var files []string

	paginator := s3.NewListObjectsV2Paginator(b.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(b.bucketName),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %v", err)
		}

		for _, object := range page.Contents {
			files = append(files, *object.Key)
		}
	}

	return files, nil
}

func (b *B2Service) CreateBucket(ctx context.Context, bucketName string) error {
	_, err := b.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}
	return nil
}

func (b *B2Service) DeleteBucket(ctx context.Context, bucketName string) error {
	_, err := b.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %v", err)
	}
	return nil
}

func (b *B2Service) SetBucketACL(ctx context.Context, bucketName, acl string) error {
	if acl != "private" && acl != "public-read" {
		return fmt.Errorf("invalid ACL: only 'private' and 'public-read' are supported")
	}

	_, err := b.client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
		Bucket: aws.String(bucketName),
		ACL:    types.BucketCannedACL(acl),
	})
	if err != nil {
		return fmt.Errorf("failed to set bucket ACL: %v", err)
	}
	return nil
}
