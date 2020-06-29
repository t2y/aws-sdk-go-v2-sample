package s3fs

import (
	"context"
	"fmt"
	"path/filepath"
)

type S3 struct {
	uri    string
	folder string
	s3uri  *S3URI

	client *S3Client
}

func (s *S3) Name() string {
	return "s3fs(" + s.s3uri.String() + ")"
}

func (s *S3) SetUriAndValidate(uri string) error {
	s.uri = uri
	s.s3uri, _ = ParseS3URI(uri)
	s.folder = s.s3uri.Key
	if !s.s3uri.IsFolder {
		s.folder = filepath.Dir(s.s3uri.Key)
	}
	return s.client.HeadObject(s.s3uri.Bucket, s.s3uri.Key)
}

func (s *S3) ListFolder() {
	ctx := context.Background()
	ch := s.client.ListObjects(ctx, s.s3uri.Bucket, s.s3uri.Key)
	for key := range ch {
		fmt.Printf("key: %s\n", key)
	}
}

func NewS3Fs(client *S3Client) *S3 {
	if client == nil {
		client = MustNewS3Client()
	}
	return &S3{
		client: client,
	}
}
