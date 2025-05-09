// Package secure provides methods for hashing and encryption.
package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrCipherTooShort reports cipher key shorter than 32 bits.
	ErrCipherTooShort = errors.New("cipher too short, must use 32-bit")
)

// HashPassword hashes a raw password string with bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt: %w", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword compares raw password from request with the hashed version.
func VerifyPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// Encrypt encrypts a raw string using the provided key and returns the encrypted version.
func Encrypt(plainText, keyRaw string) (string, error) {
	key := []byte(keyRaw)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// generate a random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// encrypt and authenticate the plaintext
	// the nonce is prepended to the ciphertext
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a string that was encrypted using Encrypt.
func Decrypt(encryptedText, keyRaw string) (string, error) {
	key := []byte(keyRaw)

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()

	// ensure the ciphertext is large enough to contain a nonce
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("key for nonce: %w", ErrCipherTooShort)
	}

	// extract nonce from the beginning of the ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
