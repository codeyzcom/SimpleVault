package web

import "SimpleVault/internal/vault"

type RecordInput struct {
	Title string
	Type  vault.RecordType

	Note       *vault.NoteData
	File       *vault.FileData
	Credential *vault.CredentialData
}
