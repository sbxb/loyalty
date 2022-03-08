package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

const (
	secretKey    = "my-super-secret-key"
	signatureKey = "another-super-secret-key"
)

const (
	signBytes = 32
	signChars = signBytes * 2
)

func encryptString(str string, key string) (string, error) {
	// get the key's hash (256 bits)
	keyBytes := sha256.Sum256([]byte(key))

	// use AES-256
	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(str), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptString(str string, key string) (string, error) {
	// get the key's hash (256 bits)
	keyBytes := sha256.Sum256([]byte(key))

	// use AES-256
	block, err := aes.NewCipher(keyBytes[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encrypted, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", errors.New("encrypted string is too short to be processed")
	}

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	result, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func makeSignature(str string, key string) string {
	secret := md5.Sum([]byte(key))
	h := hmac.New(sha256.New, secret[:])
	h.Write([]byte(str))
	signature := h.Sum(nil)

	return hex.EncodeToString(signature)
}

func GetSignedString(str string, key string) string {
	return makeSignature(str, key) + str
}

func GetStringPart(str string) string {
	if len(str) <= 64 {
		return ""
	}
	return str[signChars:]
}

func CheckSignedString(str string, key string) bool {
	if len(str) <= 64 {
		return false
	}

	value := str[signChars:]

	sign, err := hex.DecodeString(str[:signChars])
	if err != nil {
		return false
	}

	newSign, err := hex.DecodeString(makeSignature(value, key))
	if err != nil {
		return false
	}

	return hmac.Equal(newSign, sign)
}
