package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(text []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error when create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error when create GCM cipher: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error when create nonce: %w", err)
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}

func Decrypt(encryptedText []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error when create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error when create GCM cipher: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedText) < nonceSize {
		return nil, fmt.Errorf("error: encryptedText is %d but nonceSize is %d", len(encryptedText), nonceSize)
	}

	nonce, encryptedText := encryptedText[:nonceSize], encryptedText[nonceSize:]
	text, err := gcm.Open(nil, nonce, encryptedText, nil)
	if err != nil {
		return nil, fmt.Errorf("error when decrypt text: %w", err)
	}

	return text, nil
}

// turns any given key to a 32-len byte slice
func PopulateKey(key string) []byte {
	byteKey := []byte(key)
	if len(byteKey) > 32 {
		return byteKey[:32]
	}

	for i := len(byteKey); i < 32; i++ {
		byteKey = append(byteKey, byte(i))
	}

	return byteKey
}
