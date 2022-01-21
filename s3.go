package s3

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// This assumes that the AWS CLI is setup with accessible creds
// either in ~/.aws/credentials or in ENV vars
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

type Client struct {
	session *session.Session
}

func NewClient() (*Client, error) {
	client := &Client{}
	session, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String("us-east-2")},
	})
	if err != nil {
		return &Client{}, err
	}

	client.session = session
	return client, err
}

func (s *Client) Upload(bucketName string, filePath string, fileName string) error {
	uploader := s3manager.NewUploader(s.session)

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filePath, err)
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   f,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Client) Download(bucketName string, fileName string, downloadDirectory string) (string, error) {
	filePath := downloadDirectory + "/" + fileName
	downloader := s3manager.NewDownloader(s.session)

	f, err := os.Create(downloadDirectory + "/" + fileName)
	if err != nil {
		return "", fmt.Errorf("Failed to create file %v for downloading, %v", filePath, err)
	}

	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v, error: %v", fileName, err)
	}

	return filePath, nil
}

//GetFileListFromBucket - Query AWS S3 bucket for file list
//NOTE: returns max 1000 items
func (s *Client) ListBucket(bucketName string) ([]string, error) {
	var files []string
	svc := s3.New(s.session)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		return []string{}, fmt.Errorf("Error running svc.ListObjects for bucket %v - Error: %v", bucketName, err)
	}

	for _, key := range resp.Contents {
		files = append(files, *key.Key)
	}

	return files, nil
}

//DoesFileExistInBucket - Query AWS S3 bucket to see if it contains a given filename
func (s *Client) FileExists(bucketName string, fileName string) (bool, error) {
	svc := s3.New(s.session)
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(fileName),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		return false, fmt.Errorf("Error running svc.ListObjects for bucket %v - Error: %v", bucketName, err)
	}

	for _, key := range resp.Contents {
		fmt.Printf("key: %v\n", *key.Key)
		if strings.Contains(*key.Key, fileName) {
			return true, nil
		}
	}

	return false, nil
}