package password

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func GetPasswordHashAndSalt(passwd []byte) (hash, salt []byte, err error) {
	salt, err = GetSalt(20)
	if err != nil {
		return nil, nil, err
	}
	hash = GetPasswdHash([]byte(passwd), salt)
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

func GetPasswdHash(passwd, salt []byte) []byte {
	return pbkdf2.Key(passwd, salt, 4096, sha256.Size, sha256.New)
}

func ValidatePasswd(passwd, hash, salt []byte) bool {
	return bytes.Equal(GetPasswdHash(passwd, salt), hash)
}
