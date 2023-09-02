package uploader

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"room-booking/app/config"
	"room-booking/pb"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ImageUploader interface {
	upload(image bytes.Buffer, imageInfo *pb.ImageInfo) (string, error)
}

type S3Uploader struct{}

const (
	timeout = 2 * time.Second
)

// return new instance of S3Uploader
func NewS3Uploader() S3Uploader {
	return S3Uploader{}
}

func (u S3Uploader) upload(image bytes.Buffer, imageInfo *pb.ImageInfo) (string, error) {
	bucket := config.Get("AWS_S3_BUCKET")

	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(config.Get("AWS_S3_REGION")),
		Credentials: credentials.NewStaticCredentials(config.Get("AWS_S3_KEY"), config.Get("AWS_S3_SECRET"), ""),
	}))
	svc := s3.New(sess)

	ctx := context.Background()
	var cancelFn func()
	ctx, cancelFn = context.WithTimeout(ctx, timeout)
	if cancelFn != nil {
		defer cancelFn()
	}

	objectKey := fmt.Sprintf("%d-%s", time.Now().Unix(), imageInfo.GetImageName())
	_, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(image.Bytes()),
		ContentType:   aws.String(imageInfo.GetImageType()),
		ContentLength: aws.Int64(int64(image.Len())),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}

		return "", err
	}

	fmt.Printf("successfully uploaded image to S3 %s/%s\n", bucket, imageInfo.GetImageName())

	return getObjectPublicURL(objectKey), nil
}

func getObjectPublicURL(key string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.Get("AWS_S3_BUCKET"), key)
}
