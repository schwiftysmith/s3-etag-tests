package checksum

import (
	"os"

	"github.com/peak/s3hash"
)

func CalculateChecksum(path string, chunkSizeInMB int64) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	chunkSize := chunkSizeInMB

	return s3hash.Calculate(f, chunkSize)
}
