package test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"s3-etag-tests/pkg/awsclient"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

var currentDir, _ = os.Getwd()
var tempPath = fmt.Sprintf("%s/tmp", currentDir)

const defaultRegion = "us-east-1"

func createTempPath(t *testing.T, path string) {
	err := os.Mkdir(path, 0755)
	assert.NoError(t, err)
}

func removeTempPath(t *testing.T, path string) {
	if _, err := os.Stat(path); err != nil {
		pathDoesNotExist := os.IsNotExist(err)
		assert.True(t, pathDoesNotExist)
		return
	}

	err := os.RemoveAll(path)
	assert.NoError(t, err)
}

func createAWSClient(t *testing.T, ctx context.Context, endpoint string, creds awsclient.ClientCredentials) awsclient.S3Client {
	client, err := awsclient.NewS3Client(ctx, endpoint, defaultRegion, creds)
	assert.NoError(t, err)

	return client
}

func createTestBucket(t *testing.T, ctx context.Context, client awsclient.S3Client, bucketName string) {
	err := client.CreateBucket(ctx, bucketName)
	assert.NoError(t, err)
}

func createTestObjectInBucket(
	t *testing.T,
	ctx context.Context,
	client awsclient.S3Client,
	bucketName string,
	filePath string,
	objectPath string,
	partSize int64) {

	err := client.PutObject(ctx, bucketName, filePath, objectPath, partSize)
	assert.NoError(t, err)
	verifyUploadedObjectEqualsLocalFile(t, ctx, client, bucketName, objectPath, filePath, tempPath, partSize)
}

func verifyUploadedObjectEqualsLocalFile(
	t *testing.T,
	ctx context.Context,
	client awsclient.S3Client,
	bucketName string,
	objectPath string,
	localFilePath string,
	tempDownloadPath string,
	partSize int64) {

	createTempPath(t, tempDownloadPath)

	targetLocation := fmt.Sprintf("%s/%s", tempDownloadPath, filepath.Base(objectPath))

	err := client.GetObject(ctx, bucketName, objectPath, targetLocation, partSize)
	assert.NoError(t, err)

	verifyFilesAreEqual(t, localFilePath, targetLocation)

	removeTempPath(t, tempPath)
}

func verifyFilesAreEqual(t *testing.T, filePath1 string, filePath2 string) {
	file1, err := ioutil.ReadFile(filePath1)
	assert.NoError(t, err)

	file2, err := ioutil.ReadFile(filePath2)
	assert.NoError(t, err)

	areEqualFiles := bytes.Equal(file1, file2)
	assert.True(t, areEqualFiles)
}

func createFileWithSize(t *testing.T, largeFilePath string, fileSize int64) {
	largeFile, err := os.OpenFile(largeFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	assert.NoError(t, err)
	err = syscall.Fallocate(int(largeFile.Fd()), 0, 1024, fileSize)
	assert.NoError(t, err)
	err = largeFile.Close()
	assert.NoError(t, err)
}

func cleanupBucket(t *testing.T, ctx context.Context, client awsclient.S3Client, bucketName string) {
	err := client.DeleteBucket(ctx, bucketName)
	assert.NoError(t, err)
}

func cleanupFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	assert.NoError(t, err)
}
