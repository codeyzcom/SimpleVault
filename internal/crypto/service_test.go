package crypto

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCryptoService(t *testing.T) {
	srv := NewCryptoService()
	require.NotNil(t, srv)
}

func TestCryptoService_DeriveKey_Deterministic(t *testing.T) {
	password := "super-secret-password"
	salt := []byte("unique-salt")

	srv := NewCryptoService()

	key1 := srv.DeriveKey(password, salt)

	key2 := srv.DeriveKey(password, salt)

	require.Len(t, key1, 32)
	require.Len(t, key2, 32)

	assert.True(
		t,
		bytes.Equal(key1, key2),
		"same password and salt must produce identical keys",
	)
}

func TestCryptoService_DeriveKey_DifferentPassword(t *testing.T) {
	passwordOne := "super-secret-password-1"
	passwordTwo := "super-secret-password-2"
	salt := []byte("unique-salt")

	srv := NewCryptoService()

	key1 := srv.DeriveKey(passwordOne, salt)
	key2 := srv.DeriveKey(passwordTwo, salt)

	require.Len(t, key1, 32)
	require.Len(t, key2, 32)

	assert.False(
		t,
		bytes.Equal(key1, key2),
		"different passwords must produce different keys",
	)
}

func TestCryptoService_DeriveKey_DifferentSalt(t *testing.T) {
	password := "super-secret-password"

	srv := NewCryptoService()

	key1 := srv.DeriveKey(password, []byte("salt-1"))
	key2 := srv.DeriveKey(password, []byte("salt-2"))

	require.Len(t, key1, 32)
	require.Len(t, key2, 32)

	assert.False(
		t,
		bytes.Equal(key1, key2),
		"different salts must produce different keys",
	)
}
