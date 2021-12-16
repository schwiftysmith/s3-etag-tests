package awsclient

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ClientCredentials contains the credentials needed to authorize
type ClientCredentials struct {
	AccessKey string
	SecretKey string
	Session   string
}

type S3Client interface {
	// CreateBucket creates a new bucket
	CreateBucket(ctx context.Context, bucketName string) error

	// DeleteBucket deletes an existing bucket even if not empty
	DeleteBucket(ctx context.Context, bucketName string) error

	// HeadObject retrieve metadata from an object
	HeadObject(ctx context.Context, bucketName string, objectPath string) (*s3.HeadObjectOutput, error)

	// GetObject download files from bucket
	GetObject(ctx context.Context, bucketName string, objectPath string, targetFilePath string, partSize int64) error

	// PutObject upload files to bucket
	PutObject(ctx context.Context, bucketName string, localFilePath string, ObjectPath string, partSize int64) error
}

type awsS3Client struct {
	s3Client *s3.Client
}

func NewS3Client(
	ctx context.Context,
	s3Endpoint string,
	defaultRegion string,
	creds ClientCredentials) (S3Client, error) {

	s3Cfg, err := createConfig(ctx, s3Endpoint, creds, defaultRegion)

	if err != nil {
		return &awsS3Client{}, err
	}

	client := awsS3Client{
		s3Client: s3.NewFromConfig(*s3Cfg),
	}

	return &client, nil
}

func createConfig(ctx context.Context, endpoint string, creds ClientCredentials, defaultRegion string) (*aws.Config, error) {
	credsProvider := config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(creds.AccessKey, creds.SecretKey, creds.Session))
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint,
			SigningRegion:     defaultRegion,
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(defaultRegion), credsProvider, config.WithEndpointResolver(staticResolver))
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (ac *awsS3Client) CreateBucket(ctx context.Context, bucketName string) error {
	input := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	_, err := ac.s3Client.CreateBucket(ctx, input)
	return err
}

func (ac *awsS3Client) DeleteBucket(ctx context.Context, bucketName string) error {
	input := &s3.DeleteBucketInput{
		Bucket: &bucketName,
	}

	_, err := ac.s3Client.DeleteBucket(ctx, input)
	return err
}

func (ac *awsS3Client) HeadObject(ctx context.Context, bucketName string, fileName string) (*s3.HeadObjectOutput, error) {
	input := &s3.HeadObjectInput{
		Bucket: &bucketName,
		Key:    &fileName,
	}

	fileInfo, err := ac.s3Client.HeadObject(ctx, input)
	return fileInfo, err
}

func (ac *awsS3Client) GetObject(ctx context.Context, bucketName string, objectPath string, targetFilePath string, partSize int64) error {
	input := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectPath,
	}

	fileWriter, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	downloader := manager.NewDownloader(ac.s3Client)
	_, err = downloader.Download(ctx, fileWriter, input, func(u *manager.Downloader) {
		u.PartSize = partSize
	})
	return err
}

func (ac *awsS3Client) PutObject(ctx context.Context, bucketName string, localFilePath string, ObjectPath string, partSize int64) error {
	// do we really need to open the file
	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &ObjectPath,
		Body:   file,
	}

	uploader := manager.NewUploader(ac.s3Client, func(u *manager.Uploader) {
		u.PartSize = partSize
	})
	_, err = uploader.Upload(ctx, input)
	return err
}
