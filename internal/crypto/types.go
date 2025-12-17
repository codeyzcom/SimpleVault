package crypto

type Key []byte

type Keyset struct {
	Salt        []byte
	MKEncrypted []byte // Master Key encrypted by PK
	VKEncrypted []byte // Vault Key encrypted by MK
}
