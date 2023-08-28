package utils

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"time"
)

// UploadFileToS3Bucket uploads a file to an S3 bucket.

func UploadFileToS3Bucket(conf Config, key string, folder string, fileBytes []byte) error {

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.AWSRegion),
		Credentials: credentials.NewStaticCredentials(conf.AWSAccessKeyID, conf.AWSSecretAccessKey, ""),
	})

	// Create S3 service client
	svc := s3.New(sess)

	// Create the full S3 key including the folder
	fullKey := folder + "/" + key

	// Upload the file's bytes to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        &conf.AWSBucketName,
		Key:           aws.String(fullKey),
		Body:          bytes.NewReader(fileBytes),
		ContentLength: aws.Int64(int64(len(fileBytes))),
		ContentType:   aws.String("application/pdf"),
	})

	if err != nil {
		log.Error().Err(err).Msg("Could not upload file to S3")
		return err
	}
	log.Info().Msg("Successfully uploaded file to S3")
	return nil
}

// GeneratePresignedURL generates a presigned URL to download a file from S3.
func GeneratePresignedURL(conf Config, folder, key string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.AWSRegion),
		Credentials: credentials.NewStaticCredentials(conf.AWSAccessKeyID, conf.AWSSecretAccessKey, ""),
	})

	// Create S3 service client
	svc := s3.New(sess)

	// Folder path without year and week number
	folderPath := fmt.Sprintf("%s/%s", folder, key)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(conf.AWSBucketName),
		Key:    aws.String(folderPath),
	})

	presignedURL, err := req.Presign(15 * time.Minute) // URL valid for 15 minutes

	if err != nil {
		return "", err
	}

	return presignedURL, nil
}
