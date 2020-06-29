package s3fs

import (
	"fmt"
	"net/url"
	"strings"
)

type S3URI struct {
	URI      string
	Bucket   string
	Key      string
	IsFolder bool
}

func (u *S3URI) String() string {
	var objectType string
	if u.IsFolder {
		objectType = "folder"
	} else {
		objectType = "object"
	}

	tmpl := "bucket: %s, key: %s, type: %s"
	return fmt.Sprintf(tmpl, u.Bucket, u.Key, objectType)
}

func ParseS3URI(uri string) (s3uri *S3URI, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	key := u.Path[1:]
	s3uri = &S3URI{
		URI:      uri,
		Bucket:   u.Host,
		Key:      key,
		IsFolder: strings.HasSuffix(key, "/"),
	}
	return
}
