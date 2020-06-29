package main

import (
	"flag"
	"log"

	"github.com/t2y/aws-sdk-go-v2-sample/client/s3fs"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("./main s3://path/to/folder")
	}
	uri := args[0]

	s3 := s3fs.NewS3Fs(nil)
	err := s3.SetUriAndValidate(uri)
	if err != nil {
		log.Fatalf("%s is not exist", uri)
	}
	s3.ListFolder()
}
