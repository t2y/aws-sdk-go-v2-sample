package s3fs_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/t2y/aws-sdk-go-v2-sample/client/s3fs"
	"github.com/t2y/aws-sdk-go-v2-sample/client/s3mock"
)

const (
	testFolder = "s3://bucket/filter/"
)

func newMockS3(t *testing.T) (*s3mock.MockClientAPI, *s3mock.MockUploaderAPI) {
	ctrl := gomock.NewController(t)
	svc := s3mock.NewMockClientAPI(ctrl)
	uploader := s3mock.NewMockUploaderAPI(ctrl)
	return svc, uploader
}

func setupHeadObject(svc *s3mock.MockClientAPI) {
	output := &s3.HeadObjectOutput{}
	svc.EXPECT().
		HeadObjectRequest(gomock.Any()).
		Return(mockHeadObjectRequest(output, nil))
}

func setupForSetUriAndValidate(t *testing.T) *s3fs.S3 {
	svc, uploader := newMockS3(t)
	setupHeadObject(svc)
	c := s3fs.NewS3ClientForMock(svc, uploader)
	return s3fs.NewS3Fs(c)
}

func setupForListFolder(t *testing.T) *s3fs.S3 {
	svc, uploader := newMockS3(t)
	setupHeadObject(svc)

	objects := []s3.ListObjectsV2Output{
		{
			Contents: []s3.Object{
				{Key: aws.String("path/to/object1")},
				{Key: aws.String("path/to/object2")},
				{Key: aws.String("path/to/object3")},
			},
		},
	}
	svc.EXPECT().
		ListObjectsV2Request(gomock.Any()).
		Return(mockListObjectsRequest(objects, nil))

	c := s3fs.NewS3ClientForMock(svc, uploader)
	return s3fs.NewS3Fs(c)
}

func TestSetUriAndValidate(t *testing.T) {
	var data = []struct {
		name string
		uri  string
	}{
		{
			name: "path to a folder",
			uri:  "s3://bucket/folder/",
		},

		{
			name: "path to an object",
			uri:  "s3://bucket/folder/test.json",
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			s3 := setupForSetUriAndValidate(t)
			actual := s3.SetUriAndValidate(tt.uri)
			if actual != nil {
				t.Errorf("expects no error, but got %v", actual)
			}
		})
	}
}

func TestListFolder(t *testing.T) {
	var data = []struct {
		name string
		uri  string
	}{
		{
			name: "path to a folder",
			uri:  "s3://bucket/folder/",
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			s3 := setupForListFolder(t)
			err := s3.SetUriAndValidate(tt.uri)
			if err != nil {
				t.Fatalf("expects no error, but got %v", err)
			}
			s3.ListFolder()
		})
	}
}

// TODO: these mock requests are a workaround.
// see also: https://github.com/aws/aws-sdk-go-v2/issues/70
func mockHeadObjectRequest(
	output *s3.HeadObjectOutput, err error,
) s3.HeadObjectRequest {
	req := &aws.Request{
		HTTPRequest: &http.Request{},
		Retryer:     aws.NoOpRetryer{},
		Data:        output,
		Error:       err,
	}
	return s3.HeadObjectRequest{
		Request: req,
	}
}

func mockGetObjectRequest(
	output *s3.GetObjectOutput, err error,
) s3.GetObjectRequest {
	req := &aws.Request{
		HTTPRequest: &http.Request{},
		Retryer:     aws.NoOpRetryer{},
		Data:        output,
		Error:       err,
	}
	return s3.GetObjectRequest{
		Request: req,
	}
}

func mockPutObjectRequest(
	output *s3.PutObjectOutput, err error,
) s3.PutObjectRequest {
	req := &aws.Request{
		HTTPRequest: &http.Request{},
		Retryer:     aws.NoOpRetryer{},
		Data:        output,
		Error:       err,
	}
	return s3.PutObjectRequest{
		Request: req,
	}
}

func mockListObjectsRequest(
	objects []s3.ListObjectsV2Output, err error,
) s3.ListObjectsV2Request {
	i := 0
	req := s3.ListObjectsV2Request{
		Copy: func(v *s3.ListObjectsV2Input) s3.ListObjectsV2Request {
			r := s3.ListObjectsV2Request{
				Request: &aws.Request{
					HTTPRequest: &http.Request{},
					Retryer:     aws.NoOpRetryer{},
					Operation:   &aws.Operation{},
					Error:       err,
				},
			}
			r.Handlers.Send.PushBack(func(r *aws.Request) {
				obj := objects[i]
				r.Data = &obj
				i++
			})
			return r
		},
	}
	return req
}
