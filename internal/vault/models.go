package vault

import "time"

type RecordType string

const (
	RecordNote       RecordType = "note"
	RecordFile       RecordType = "file"
	RecordCredential RecordType = "credential"
)

type Meta struct {
	Salt     []byte `json:"salt"`
	Verifier []byte `json:"verifier"`
}

type Vault struct {
	Records []*Record `json:"records"`
}

type Record struct {
	ID        string `json:"id"`
	Title     string
	Type      RecordType `json:"type"`
	CreatedAt time.Time  `json:"created_at"`

	Note       *NoteData
	File       *FileData
	Credential *CredentialData
}

type (
	NoteData struct {
		Text string
	}

	FileData struct {
		Filename string
		Data     []byte
	}

	CredentialData struct {
		Website  string
		Username string
		Password string
		Email    string
		Phone    string
		Note     string
	}
)
