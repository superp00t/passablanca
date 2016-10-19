package cryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"

	"io"
)

func DeriveKey(password string) []byte {
	bytes := sha256.Sum256([]byte(password))
	return bytes[:32]
}

func Encrypt(password string, plaintext []byte) []byte {
	key := DeriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func Decrypt(password string, ciphertext []byte) []byte {
	key := DeriveKey(password)

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("cipherText too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}
