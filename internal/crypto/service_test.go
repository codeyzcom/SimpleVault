package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
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

func TestCryptoService_VerifyKey(t *testing.T) {
	svc := &CryptoService{}

	t.Run("valid key and verifier", func(t *testing.T) {
		key := []byte("super-secret-password")
		hash := sha256.Sum256(key)

		ok := svc.VerifyKey(key, hash[:])

		assert.True(t, ok)
	})

	t.Run("invalid key", func(t *testing.T) {
		key := []byte("super-secret-password")
		wrongKey := []byte("another-password")

		hash := sha256.Sum256(key)

		ok := svc.VerifyKey(wrongKey, hash[:])

		assert.False(t, ok)
	})

	t.Run("empty key", func(t *testing.T) {
		key := []byte{}
		hash := sha256.Sum256(key)

		ok := svc.VerifyKey([]byte{}, hash[:])

		assert.True(t, ok)
	})

	t.Run("empty verifier", func(t *testing.T) {
		key := []byte("secret")

		ok := svc.VerifyKey(key, []byte{})

		assert.False(t, ok)
	})

	t.Run("corrupted verifier", func(t *testing.T) {
		key := []byte("secret")
		hash := sha256.Sum256(key)

		verifier := hash[:]
		verifier[0] ^= 0xff

		ok := svc.VerifyKey(key, verifier)

		assert.False(t, ok)
	})
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	srv := NewCryptoService()

	key := make([]byte, 32)
	_, _ = rand.Read(key)

	plaintext := []byte("top secret payload")

	encrypted, err := srv.Encrypt(key, plaintext)
	require.NoError(t, err)
	require.NotEmpty(t, encrypted)

	decrypted, err := srv.Decrypt(key, encrypted)
	require.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}
