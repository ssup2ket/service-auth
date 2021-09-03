package password

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

func GetPasswordHashAndSalt(passwd string) (hash, salt []byte, err error) {
	salt, err = GetSalt(20)
	if err != nil {
		return nil, nil, err
	}
	hash = GetPasswdHash(passwd, salt)
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

func GetPasswdHash(passwd string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passwd), salt, 4096, sha256.Size, sha256.New)
}

func ValidatePasswd(passwd string, hash, salt []byte) bool {
	return bytes.Equal(GetPasswdHash(passwd, salt), hash)
}
