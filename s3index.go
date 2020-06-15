// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

// +build ignore

// This tool uploads a file containing branch names and commit shas to S3 to trigger CI/CD.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	client := s3.New(session.Must(session.NewSession()), &aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})

	var buf bytes.Buffer
	buf.WriteString(os.Getenv("CIRCLE_BRANCH") + "\n")
	buf.WriteString(os.Getenv("CIRCLE_SHA1") + "\n")
	buf.WriteString(os.Getenv("CIRCLE_USERNAME") + "\n")

	uploader := s3manager.NewUploaderWithClient(client)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("datadog-reliability-env"),
		Key:    aws.String("go/index.txt"),
		Body:   &buf,
	})
	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
	}
	fmt.Printf("index.txt uploaded to, %s\n", aws.StringValue(&result.Location))
}
