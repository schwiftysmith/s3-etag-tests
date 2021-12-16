package checksum

import (
	"os"

	"github.com/peak/s3hash"
)

func CalculateChecksum(path string, chunkSize int64) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return s3hash.Calculate(f, chunkSize)
}
