package test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"s3-etag-tests/pkg/awsclient"
	"s3-etag-tests/pkg/checksum"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s3Endpoint = flag.String("s3Endpoint", "http://localhost:8333", "s3 server url")
var accessKey = flag.String("accessKey", "accessKey", "s3 access key")
var secretAccessKey = flag.String("secretAccessKey", "secretAccessKey", "s3 secret key")

const (
	ChunkSizeInByte     = 1024 * 1024 * 16
	LargeFileSizeInByte = 1024 * 1024 * 50 // > chunk size
)

var workingDir, _ = os.Getwd()
var testFileName = "test-object.txt"
var bucketName = "my-bucket"
var keyPath = ""
var objectPath = fmt.Sprintf("%s%s", keyPath, testFileName)
var dataPath = fmt.Sprintf("%s/data/%s", workingDir, testFileName)

var awsS3cfg = awsclient.ClientCredentials{
	AccessKey: *accessKey,
	SecretKey: *secretAccessKey,
	Session: "",
}

func TestEtagOfSmallFile(t *testing.T) {
	// prepare
	ctx := context.Background()
	client := createAWSClient(t, ctx, *s3Endpoint, awsS3cfg)

	createTestBucket(t, ctx, client, bucketName)
	createTestObjectInBucket(t, ctx, client, bucketName, dataPath, objectPath, ChunkSizeInByte)

	expectedEtag, err := checksum.CalculateChecksum(dataPath, int64(ChunkSizeInByte))
	assert.NoError(t, err)

	// test
	fileInfo, err := client.HeadObject(ctx, bucketName, objectPath)
	actualEtag := strings.Trim(*fileInfo.ETag, "\"")
	assert.NoError(t, err)
	assert.Equal(t, expectedEtag, actualEtag)

	// cleanup
	cleanupBucket(t, ctx, client, bucketName)
}

func TestEtagOfLargeFile(t *testing.T) {
	// prepare
	ctx := context.Background()
	client := createAWSClient(t, ctx, *s3Endpoint, awsS3cfg)

	// prepare large file
	largeFileName := "large-blobb.txt"
	largeFilePath := fmt.Sprintf("%s/%s", workingDir, largeFileName)
	largeObjectPath := fmt.Sprintf("%s%s", keyPath, largeFileName)
	createFileWithSize(t, largeFilePath, LargeFileSizeInByte)

	createTestBucket(t, ctx, client, bucketName)
	createTestObjectInBucket(t, ctx, client, bucketName, largeFilePath, largeObjectPath, ChunkSizeInByte)

	expectedEtag, err := checksum.CalculateChecksum(largeFilePath, int64(ChunkSizeInByte))
	assert.NoError(t, err)

	// test
	fileInfo, err := client.HeadObject(ctx, bucketName, largeObjectPath)
	actualEtag := strings.Trim(*fileInfo.ETag, "\"")
	assert.NoError(t, err)
	assert.Equal(t, expectedEtag, actualEtag) //This fails while using SeaweedFS

	// cleanup
	cleanupBucket(t, ctx, client, bucketName)
	cleanupFile(t, largeFilePath)
}
