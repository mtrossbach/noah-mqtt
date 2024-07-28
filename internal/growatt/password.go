package growatt

import (
	"crypto/md5"
	"encoding/hex"
)

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	digest := hex.EncodeToString(hash[:])

	for i := 0; i < len(digest); i = i + 2 {
		if digest[i] == '0' {
			digest = digest[:i] + "c" + digest[i+1:]
		}
	}

	return digest
}
