package hashing

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func GetStrHashAndSalt(str string) (hash, salt []byte, err error) {
	salt, err = GetSalt(20)
	if err != nil {
		return nil, nil, err
	}
	hash = GetStrHash(str, salt)
	return
}

func GetSalt(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GetStrHash(passwd string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passwd), salt, 4096, sha256.Size, sha256.New)
}

func ValidateStr(passwd string, hash, salt []byte) bool {
	return bytes.Equal(GetStrHash(passwd, salt), hash)
}
