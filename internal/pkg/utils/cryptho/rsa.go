package cryptho

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

var hashSeed = sha256.New()

// Function to generate RSA key pair and save to files
func saveRSAKeysToFile(privateKeyPath, publicKeyPath string, bits int) error {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("error generating rsa private key: %v", err)
	}

	// Save the private key in PEM format
	privateFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("error creating private key file: %v", err)
	}
	defer privateFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	err = pem.Encode(privateFile, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("error encoding private key: %v", err)
	}

	// Save the public key in PEM format
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("error marshalling public key: %v", err)
	}

	publicFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("error creating public key file: %v", err)
	}
	defer publicFile.Close()

	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	err = pem.Encode(publicFile, publicKeyPEM)
	if err != nil {
		return fmt.Errorf("error encoding public key: %v", err)
	}

	return nil
}

// Function to load the RSA private key from a PEM file
func LoadRSAPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v, filepath: %v", err, filePath)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	return privateKey, nil
}

// Function to load the RSA public key from a PEM file
func LoadRSAPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %v, filepath: %v", err, filePath)
	}

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	publicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return publicKey, nil
}

// RSA encryption using a public key, with string input and base64-encoded string output
func RsaEncrypt(publicKey *rsa.PublicKey, message string) (string, error) {
	// Encrypt using OAEP with SHA-256
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(message), nil)
	if err != nil {
		return "", err
	}

	// Encode the ciphertext to base64 string
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return encodedCiphertext, nil
}

// RSA decryption using a private key, with base64-encoded string input and string output
func RsaDecrypt(privateKey *rsa.PrivateKey, ciphertext string) (string, error) {
	// Decode the base64 encoded ciphertext
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Decrypt using OAEP with SHA-256
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, decodedCiphertext, nil)
	if err != nil {
		return "", err
	}

	// Convert the plaintext bytes back to a string
	return string(plaintext), nil
}
