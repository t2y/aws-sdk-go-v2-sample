package client_test

//go:generate mockgen -package sqsmock -destination sqsmock/sqs.go github.com/aws/aws-sdk-go-v2/service/sqs/sqsiface ClientAPI
