package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFull(t *testing.T) {
	message := "abcdefghijklmnopqrstuvwxyz"
	cryptoKey := "secretCrypto"
	signKey := "secretSign"

	encryptedMessage, err := encryptString(message, cryptoKey)
	require.NoError(t, err)

	signedEncMessage := GetSignedString(encryptedMessage, signKey)
	t.Log(signedEncMessage)
	assert.True(t, CheckSignedString(signedEncMessage, signKey))

	encryptedMessage2 := GetStringPart(signedEncMessage)

	decryptedMessage, err := decryptString(encryptedMessage2, cryptoKey)
	require.NoError(t, err)

	assert.Equal(t, decryptedMessage, message)
}

func TestEncryptDecrypt(t *testing.T) {
	key := "secret"
	value := "SomeTextInfoToEncrypt"

	encrypted, err := encryptString(value, key)
	require.NoError(t, err)

	decrypted, err := decryptString(encrypted, key)
	require.NoError(t, err)

	assert.Equal(t, decrypted, value)
}

func TestSignCheck(t *testing.T) {
	key := "sigsecret"
	value := "SomeTextInfoToSign"

	signedString := GetSignedString(value, key)

	assert.Equal(t, CheckSignedString(signedString, key), true)
}
