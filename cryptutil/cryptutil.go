package cryptutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	rd "math/rand"
	"time"
)

type CryptStore struct {
	Nonce      string
	Ciphertext string
}

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

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	cs := CryptStore{
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}

	dat, _ := json.Marshal(cs)

	return dat
}

func Decrypt(password string, ctext []byte) ([]byte, error) {
	var cs CryptStore

	err := json.Unmarshal(ctext, &cs)
	if err != nil {
		panic(err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(cs.Ciphertext)
	if err != nil {
		panic(err)
	}
	nonce, err := base64.StdEncoding.DecodeString(cs.Nonce)
	if err != nil {
		panic(err)
	}
	key := DeriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func RandomString() string {
	rd.Seed(time.Now().UnixNano())
	sl := make([]byte, 36)
	rand.Read(sl)
	return base64.StdEncoding.EncodeToString(sl)
}
