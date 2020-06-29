package client_test

//go:generate mockgen -package s3mock -destination s3mock/s3.go github.com/aws/aws-sdk-go-v2/service/s3/s3iface ClientAPI
//go:generate mockgen -package s3mock -destination s3mock/s3manager.go github.com/aws/aws-sdk-go-v2/service/s3/s3manager/s3manageriface UploaderAPI
