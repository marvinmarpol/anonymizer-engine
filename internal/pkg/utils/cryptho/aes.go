package cryptho

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AESEncrypt encrypts plain text using AES-GCM
func AESEncrypt(plainText, key string) (string, error) {
	// Convert key to byte slice
	keyBytes := []byte(key)

	// Create new AES cipher with the given key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Create GCM mode on the cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce with the required size
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)

	// Return base64 encoded string
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// AESDecrypt AES-GCM
func AESDecrypt(cipherText, key string) (string, error) {
	// Convert key to byte slice
	keyBytes := []byte(key)

	// Base64 decode the cipher text
	cipherData, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	// Create new AES cipher with the given key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Create GCM mode on the cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Check if cipher text is too short
	if len(cipherData) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Split nonce and encrypted message
	nonce, cipherTextBytes := cipherData[:nonceSize], cipherData[nonceSize:]

	// Decrypt the message
	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", err
	}

	// Return the decrypted message as a string
	return string(plainText), nil
}
