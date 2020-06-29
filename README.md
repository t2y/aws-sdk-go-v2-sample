# aws-sdk-go-v2-sample

aws-sdk-go-v2 sample code

## How to test

Generate mock code for AWS service.

### S3

```bash
$ make s3generate
go generate client/s3fs_test.go
$ tree client/s3mock
client/s3mock
├── s3.go
└── s3manager.go
```

This repository only have a test for S3 service.

* `client/s3fs/s3fs_test.go`

To run the test with s3 mock is like this.

```bash
$ make test
```

### SQS

Generating sqs mock code only, there are no tests for SQS.

```bash
$ make sqsgenerate
go generate client/sqs_test.go
$ tree client/sqsmock
client/sqsmock
└── sqs.go
```

## References

* https://github.com/aws/aws-sdk-go-v2
* [Mocking Out the AWS SDK for Go for Unit Testing](https://aws.amazon.com/jp/blogs/developer/mocking-out-then-aws-sdk-for-go-for-unit-testing/)
* [Mocking out new API? #70](https://github.com/aws/aws-sdk-go-v2/issues/70)
