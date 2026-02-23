package hash

import (
	"io"
	"os"

	"github.com/cespare/xxhash/v2"
)

// Hash64 计算给定字节数组的 xxHash64 哈希值.
func Hash64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

// HashReader 计算任何流式数据 (io.Reader) 的哈希值
func HashReader(r io.Reader) (uint64, error) {
	hasher := xxhash.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return 0, err
	}
	return hasher.Sum64(), nil
}

// HashFile 计算指定本地文件的 xxHash64 哈希值.
func HashFile(filePath string) (uint64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return HashReader(file)
}