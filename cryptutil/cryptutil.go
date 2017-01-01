package cryptutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"io"
	"strings"

	"github.com/superp00t/niceware"
)

type CryptStore struct {
	Nonce      []byte
	Ciphertext []byte
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
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}

	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(cs)

	return buf.Bytes()
}

func Decrypt(password string, ctext []byte) ([]byte, error) {
	var cs CryptStore

	err := gob.NewDecoder(bytes.NewReader(ctext)).Decode(&cs)
	if err != nil {
		return nil, err
	}

	key := DeriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, cs.Nonce, cs.Ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func RandomString() string {
	pass, _ := niceware.RandomString()
	fixedpass := strings.Replace(pass, " ", "_", -1)
	return fixedpass
}
