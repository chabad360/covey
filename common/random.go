package common

import (
	"fmt"

	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// RandomString generates a random 32 byte random string.
func RandomString() string {
	b := make([]byte, 32)

	for i := range b {
		r, err := rand.Int(rand.Reader, big.NewInt(0xffffffff))
		if err != nil {
			panic(err)
		}

		b[i] = charset[r.Int64()%int64(len(charset))]
	}

	return string(b)
}

// GenerateID takes an object, converts it to text (adds two 32 byte random strings) and returns a sha256 hash from it.
func GenerateID(item interface{}) string {
	id := sha256.Sum256([]byte(fmt.Sprintf("%v%v%v", RandomString(), item, RandomString())))
	return hex.EncodeToString(id[:])
}
