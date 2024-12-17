package cryptho

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
)

// HashType represents the type of hash (sha1, sha256, sha512)
type HashType string

const (
	SHA1   HashType = "sha1"
	SHA256 HashType = "sha256"
	SHA512 HashType = "sha512"
	MD5    HashType = "md5"
)

// hashFunctionsMap associates hash types with their respective hash functions
var hashFunctionsMap = map[HashType]func() hash.Hash{
	SHA1:   sha1.New,
	SHA256: sha256.New,
	SHA512: sha512.New,
	MD5:    md5.New,
}

// GenerateHash generates a hash for the given data and hash type
func GenerateHash(input string, hashType HashType) (string, error) {
	// Get the hash function based on the hashType
	hashFunc, exists := hashFunctionsMap[hashType]
	if !exists {
		return "", errors.New("unsupported hash type")
	}

	// Create the hash and compute the sum
	h := hashFunc()
	h.Write([]byte(input))

	// Return the result as a hex string
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
