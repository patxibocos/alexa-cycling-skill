package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/patxibocos/alexa-cycling-skill/alexa-skill-lambda/pcsscraper"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
)

var s3Client *s3.S3

func getS3Client() *s3.S3 {
	if s3Client == nil {
		sess := session.Must(session.NewSession())
		s3Client = s3.New(sess)
	}
	return s3Client
}

func HandleRequest(ctx context.Context) (string, error) {
	s3Bucket := os.Getenv("AWS_S3_BUCKET")
	s3ObjectKey := os.Getenv("AWS_S3_OBJECT_KEY")
	output, _ := getS3Client().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3ObjectKey),
	})
	cyclingData := new(pcsscraper.CyclingData)
	body, _ := ioutil.ReadAll(output.Body)
	_ = proto.Unmarshal(body, cyclingData)
	whatever := cyclingData.Riders[0].LastName
	return fmt.Sprintf("Hi there %s", whatever), nil
}

func main() {
	lambda.Start(HandleRequest)
}
