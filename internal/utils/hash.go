package utils

import (
	"crypto/md5"
	"encoding/binary"
)

func HashFunction(value string) int {
	hash := md5.Sum([]byte(value))
	return int(binary.BigEndian.Uint32(hash[:4])) % 100
}
