package crypto

import "golang.org/x/crypto/argon2"

func DerivePasswordKey(password string, salt []byte) Key {
	return argon2.IDKey(
		[]byte(password),
		salt,
		3,       // time
		64*1024, // memory (64 MB)
		4,       //threads
		32,      // key length
	)
}
