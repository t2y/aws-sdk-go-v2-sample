.PHONY: all
all: build

NOW 		:=	$(shell date --utc +'%Y-%m-%dT%TZ')
REVISION 	:=	$(shell git rev-parse --short=8 HEAD)

LDFLAGS		:=	"-X main.revision=$(REVISION) -X main.buildTime=$(NOW)"

.PHONY: build
build:
	go build -o main -ldflags $(LDFLAGS) main.go

.PHONY: modclean
modclean:
	go mod tidy

.PHONY: clean
clean:
	go clean -testcache
	rm -f main
	rm -rf client/*mock

.PHONY: s3generate
s3generate:
	go generate client/s3mock_test.go

client/s3mock/s3.go: s3generate
client/s3mock/s3manager.go: s3generate

.PHONY: sqsgenerate
sqsgenerate:
	go generate client/sqsmock_test.go

client/sqsmock/sqs.go: sqsgenerate

.PHONY: test
test: client/s3mock/s3.go
	go test -v -race -cover ./...
