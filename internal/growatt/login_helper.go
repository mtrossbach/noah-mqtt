package growatt

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
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

func timestamp() string {
	// Get the current time in milliseconds and convert it to a string
	valueOf := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Extract characters at positions 1, 3, 5, and 7
	extracted := string(valueOf[1]) + string(valueOf[3]) + string(valueOf[5]) + string(valueOf[7])

	// Parse the extracted string to an integer and take modulo 98
	parsedInt, _ := strconv.Atoi(extracted)
	parsedInt %= 98

	// Start building the final string
	result := valueOf[:11]

	// Conditionally format the parsed integer
	if parsedInt < 10 {
		result += "0" + strconv.Itoa(parsedInt)
	} else {
		result += strconv.Itoa(parsedInt)
	}

	return result
}

func ipvcpc(username string) string {
	hash := []byte(hashPassword(username + "★☆i₰₭" + fmt.Sprintf("%d", time.Now().UnixMilli())))
	hashCode := int64(binary.LittleEndian.Uint64(hash[:8]))

	// Create a UUID using the hash code as the most significant bits and a fixed number for the least significant bits.
	u := uuid.New()
	u[0] = byte(hashCode >> 56)
	u[1] = byte(hashCode >> 48)
	u[2] = byte(hashCode >> 40)
	u[3] = byte(hashCode >> 32)
	u[4] = byte(hashCode >> 24)
	u[5] = byte(hashCode >> 16)
	u[6] = byte(hashCode >> 8)
	u[7] = byte(hashCode)

	// Convert to string and return.
	return u.String()
}
