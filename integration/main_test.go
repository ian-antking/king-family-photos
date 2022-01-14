package integration

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

var testSuite *integrationTestSuite

type integrationTestSuite struct {
	suite.Suite
	environment       string
	s3Downloader      *s3manager.Downloader
	s3Uploader        *s3manager.Uploader
	s3Client          *s3.S3
	ingestBucketName  string
	displayBucketName string
}

func (s *integrationTestSuite) putImageInIngestBucket(img image.Image) string {
	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, img, nil)

	u, _ := uuid.NewUUID()
	key := u.String()

	putObjectInput := s3manager.UploadInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(s.ingestBucketName),
		Key:    aws.String(key),
	}

	_, err := s.s3Uploader.Upload(&putObjectInput)

	if nil != err {
		log.Fatalln(err.Error())
	}

	return key
}

func (s *integrationTestSuite) listItemsInBucket(bucketName string) []string {
	output, err := s.s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	})

	if nil != err {
		log.Fatalln(err.Error())
	}

	var objects []string
	for _, object := range output.Contents {
		key := object.Key
		objects = append(objects, *key)
	}

	return objects
}

func (s *integrationTestSuite) deleteObjectsInBucket(bucketName string) {
	keys := s.listItemsInBucket(bucketName)

	if 0 == len(keys) {
		return
	}

	input := s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &s3.Delete{
			Objects: []*s3.ObjectIdentifier{},
		},
	}

	for _, key := range keys {
		objectId := s3.ObjectIdentifier{
			Key: aws.String(key),
		}
		input.Delete.Objects = append(input.Delete.Objects, &objectId)
	}

	_, err := s.s3Client.DeleteObjects(&input)

	if nil != err {
		log.Fatalln(err.Error())
	}
}

func (s *integrationTestSuite) deleteObjectFromBucket(bucket, key string) {
	input := s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := s.s3Client.DeleteObject(&input)

	if nil != err {
		log.Fatalln(err.Error())
	}
}

func (s *integrationTestSuite) getImageFromBucket(bucket, key string) image.Image {
	input := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	buffer := &aws.WriteAtBuffer{}

	_, err := s.s3Downloader.Download(buffer, &input)

	if nil != err {
		log.Fatalln(err.Error())
	}

	img, _, _ := image.Decode(bytes.NewReader(buffer.Bytes()))

	return img
}

func (s *integrationTestSuite) TearDownTest() {
	s.deleteObjectsInBucket(s.displayBucketName)
	s.deleteObjectsInBucket(s.ingestBucketName)
}

func init() {
	testSuite = new(integrationTestSuite)

	flag.StringVar(&testSuite.environment, "environment", "dev", "the environment the tests are running in")

	awsSession := session.Must(session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: aws.String("eu-west-2")},
			SharedConfigState: session.SharedConfigEnable,
		},
	))

	testSuite.s3Downloader = s3manager.NewDownloader(awsSession)
	testSuite.s3Uploader = s3manager.NewUploader(awsSession)
	testSuite.s3Client = s3.New(awsSession)
	testSuite.ingestBucketName = fmt.Sprintf("king-family-photos-%s-ingest", testSuite.environment)
	testSuite.displayBucketName = fmt.Sprintf("king-family-photos-%s-display", testSuite.environment)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, testSuite)
}
