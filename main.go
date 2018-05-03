package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
)

func main() {

	c, err := GetConfig("./conf/audiofp.conf")
	if err != nil {
		fmt.Println("Config file not found")
		return
	}

	fileToUpload := os.Args[1]

	file, err := os.Open(fileToUpload)
	if err != nil {
		fmt.Println("Failed to open file", fileToUpload, err)
		os.Exit(1)
	}
	defer file.Close()

	// Initialize a client using Spaces
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(c.Bucket.AccessKey, c.Bucket.SecretKey, ""),
		Endpoint:    aws.String(c.Bucket.Endpoint),
		Region:      aws.String(c.Bucket.Region), // This is counter intuitive, but it will fail with a non-AWS region name.
	}

	newSession := session.New(s3Config)
	svc := s3manager.NewUploader(newSession)

	// Upload a file to the Space
	fmt.Println("Uploading file to S3...")
	result, err := svc.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(c.Bucket.Name),
		Key:         aws.String(filepath.Base(fileToUpload)),
		Body:        file,
		ContentType: aws.String("video/mp4"),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		fmt.Println("error ", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uploaded %s to %s\n", fileToUpload, result.Location)
}
