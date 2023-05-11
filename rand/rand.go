package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const RememberTokenBytes = 32

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nRead != n {
		return nil, fmt.Errorf("bytes: read %d bytes, expected %d", nRead, n)
	}
	return nil, nil
}

// String returns a base64 encoded string of n random bytes, or an error if there was one.
func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SessionToken is a helper function designed to generate a session token of a predetermined byte size.
func SessionToken() (string, error) {
	return String(RememberTokenBytes)
}
