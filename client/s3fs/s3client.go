package s3fs

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3iface"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager/s3manageriface"
)

type S3Client struct {
	svc      s3iface.ClientAPI
	uploader s3manageriface.UploaderAPI
}

func (c *S3Client) HeadBucket(bucket string) error {
	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}
	req := c.svc.HeadBucketRequest(input)
	_, err := req.Send(context.Background())
	return err
}

func (c *S3Client) HeadObject(bucket, key string) error {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	req := c.svc.HeadObjectRequest(input)
	_, err := req.Send(context.Background())
	return err
}

func (c *S3Client) GetObject(
	ctx context.Context, bucket, key string,
) (*s3.GetObjectResponse, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	req := c.svc.GetObjectRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		log.Printf("failed to get object: %v\n", err)
		return nil, err
	}
	return res, err
}

func (c *S3Client) PutObject(
	ctx context.Context, bucket, key string, body io.ReadSeeker,
) (*s3.PutObjectResponse, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	req := c.svc.PutObjectRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		log.Printf("failed to put object: %v\n", err)
		return nil, err
	}
	return res, err
}

func (c *S3Client) Upload(
	ctx context.Context, bucket, key string, body io.Reader,
) (*s3manager.UploadOutput, error) {
	input := &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	res, err := c.uploader.UploadWithContext(ctx, input)
	if err != nil {
		log.Printf("failed to upload: %v\n", err)
	}
	return res, err
}

func (c *S3Client) ListObjects(
	ctx context.Context, bucket string, prefix string,
) <-chan string {
	ch := make(chan string, 32)
	go func() {
		defer close(ch)

		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		}
		req := c.svc.ListObjectsV2Request(input)
		p := s3.NewListObjectsV2Paginator(req)
		for p.Next(ctx) {
			page := p.CurrentPage()
			for _, obj := range page.Contents {
				if !strings.HasSuffix(*obj.Key, "/") {
					ch <- *obj.Key
				}
			}
		}

		if err := p.Err(); err != nil {
			log.Printf("failed to list objects: %v\n", err)
		}
	}()
	return ch
}

func NewS3ClientForMock(
	svc s3iface.ClientAPI, uploader s3manageriface.UploaderAPI,
) *S3Client {
	return &S3Client{
		svc:      svc,
		uploader: uploader,
	}
}

func NewS3Client() (*S3Client, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	uploader := s3manager.NewUploader(cfg, func(u *s3manager.Uploader) {
		// Define a strategy that will buffer size in memory
		size := 32 * 1024 * 1024 // 32MiB
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(size)
	})

	c := &S3Client{
		svc:      s3.New(cfg),
		uploader: uploader,
	}
	return c, nil
}

func MustNewS3Client() *S3Client {
	c, err := NewS3Client()
	if err != nil {
		panic(err)
	}
	return c
}
