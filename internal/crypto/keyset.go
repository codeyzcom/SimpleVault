package crypto

import (
	"crypto/rand"
	"errors"
)

func CreateKeyset(password string) (*Keyset, Key, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	pk := DerivePasswordKey(password, salt)

	mk := make([]byte, 32)
	if _, err := rand.Read(mk); err != nil {
		return nil, nil, err
	}

	vk := make([]byte, 32)
	if _, err := rand.Read(vk); err != nil {
		return nil, nil, err
	}

	mkEnc, err := Encrypt(pk, mk)
	if err != nil {
		return nil, nil, err
	}

	vkEnc, err := Encrypt(mk, vk)
	if err != nil {
		return nil, nil, err
	}

	return &Keyset{
		Salt:        salt,
		MKEncrypted: mkEnc,
		VKEncrypted: vkEnc,
	}, vk, nil
}

func OpenKeyset(password string, ks *Keyset) (Key, error) {
	pk := DerivePasswordKey(password, ks.Salt)
	mk, err := Decrypt(pk, ks.MKEncrypted)
	if err != nil {
		return nil, errors.New("invalid password")
	}

	vk, err := Decrypt(mk, ks.VKEncrypted)
	if err != nil {
		return nil, err
	}
	return vk, nil
}

func CreateDataKey(valueKey Key) (Key, []byte, error) {
	dk := make([]byte, 32)
	if _, err := rand.Read(dk); err != nil {
		return nil, nil, err
	}

	dkEnc, err := Encrypt(valueKey, dk)
	if err != nil {
		return nil, nil, err
	}

	return dk, dkEnc, nil
}

func OpenDataKey(vaultKey Key, encrypted []byte) (Key, error) {
	return Decrypt(vaultKey, encrypted)
}
