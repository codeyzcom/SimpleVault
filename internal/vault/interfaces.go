package vault

type Crypto interface {
	DeriveKey(password string, salt []byte) []byte
	VerifyKey(key, verifier []byte) bool
	Encrypt(key, plaintext []byte) ([]byte, error)
	Decrypt(key, ciphertext []byte) ([]byte, error)
}

type Storage interface {
	Read(name string) ([]byte, error)
	Write(name string, data []byte) error
}
