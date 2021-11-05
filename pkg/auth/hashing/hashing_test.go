package hashing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatePasswd(t *testing.T) {
	passwd := "passwd"
	hash, salt, err := GetStrHashAndSalt(passwd)

	require.Nil(t, err, "Failed to get password hash and salt")
	require.Equal(t, true, ValidateStr(passwd, hash, salt), "Failed to vaildate password")
}
