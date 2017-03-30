package cryptoutil

import (
	"crypto/sha256"
	"fmt"
)

const (
	stretchingCount = 10
)

// GenerateHash ...
func GenerateHash(passwd string) string {
	data := []byte(passwd)
	hashed := sha256.Sum256(data)
	for i := 0; i < stretchingCount-1; i++ {
		hashed = sha256.Sum256(data)
	}
	return fmt.Sprint(hashed)
}
