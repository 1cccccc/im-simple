package utils

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
