package vault

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type VaultService struct {
	crypto  Crypto
	storage Storage

	key   []byte
	vault *Vault
	meta  *Meta
}

func NewVaultService(c Crypto, s Storage) *VaultService {
	return &VaultService{crypto: c, storage: s}
}

func (s *VaultService) Init(password string) error {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key := s.crypto.DeriveKey(password, salt)

	meta := &Meta{
		Salt:     salt,
		Verifier: sha(key),
	}

	v := &Vault{Records: []*Record{}}

	raw, _ := json.Marshal(v)
	enc, _ := s.crypto.Encrypt(key, raw)

	metaBytes, _ := json.Marshal(meta)

	if err := s.storage.Write("vault.meta", metaBytes); err != nil {
		return err
	}

	if err := s.storage.Write("vault.dat", enc); err != nil {
		return err
	}

	return nil
}

func (s *VaultService) Login(password string) error {
	metaBytes, err := s.storage.Read("vault.meta")
	if err != nil {
		return err
	}

	var meta Meta
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return err
	}

	key := s.crypto.DeriveKey(password, meta.Salt)
	if !s.crypto.VerifyKey(key, meta.Verifier) {
		return errors.New("invalid password")
	}

	data, err := s.storage.Read("vault.dat")
	if err != nil {
		return err
	}

	plain, err := s.crypto.Decrypt(key, data)
	if err != nil {
		return err
	}

	var v Vault
	if err := json.Unmarshal(plain, &v); err != nil {
		return err
	}

	s.key = key
	s.vault = &v
	s.meta = &meta
	return nil
}

func (s *VaultService) List() []*Record {
	return s.vault.Records
}

func (s *VaultService) Save() error {
	raw, err := json.Marshal(s.vault)
	if err != nil {
		return err
	}

	enc, err := s.crypto.Encrypt(s.key, raw)
	if err != nil {
		return err
	}

	return s.storage.Write("vault.dat", enc)
}

func (s VaultService) Wipe() {
	if s.key != nil {
		for i := range s.key {
			s.key[i] = 0
		}
	}

	s.key = nil
	s.vault = nil
	s.meta = nil
}

func (s *VaultService) AddNote(title, text string) error {
	if title == "" || text == "" {
		return errors.New("title and text are required")
	}

	r := &Record{
		ID:        uuid.NewString(),
		Title:     title,
		Type:      RecordNote,
		CreatedAt: time.Now(),
		Note: &NoteData{
			Text: text,
		},
	}

	s.vault.Records = append(s.vault.Records, r)
	return s.Save()
}

func (s *VaultService) AddFile(title, filename string, data []byte) error {
	if title == "" || filename == "" || len(data) == 0 {
		return errors.New("invalid file data")
	}

	if len(data) > 8*1024*1024 {
		return errors.New("file too large (max 8MB)")
	}

	r := &Record{
		ID:        uuid.NewString(),
		Title:     title,
		Type:      RecordFile,
		CreatedAt: time.Now(),
		File: &FileData{
			Filename: filename,
			Data:     data,
		},
	}

	s.vault.Records = append(s.vault.Records, r)
	return s.Save()
}

func (s *VaultService) AddCredential(title string, c CredentialData) error {
	if c.Password == "" {
		return errors.New("title and password are required")
	}

	r := &Record{
		ID:        uuid.NewString(),
		Title:     title,
		Type:      RecordCredential,
		CreatedAt: time.Now(),
		Credential: &CredentialData{
			Website:  c.Website,
			Username: c.Username,
			Password: c.Password,
			Email:    c.Email,
			Phone:    c.Phone,
			Note:     c.Note,
		},
	}

	s.vault.Records = append(s.vault.Records, r)
	return s.Save()
}

func (s *VaultService) GetRecord(id string) (*Record, error) {
	for _, r := range s.vault.Records {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, errors.New("record not found")
}

func (s *VaultService) DeleteRecord(id string) error {
	for i, r := range s.vault.Records {
		if r.ID == id {
			s.vault.Records = append(
				s.vault.Records[:i],
				s.vault.Records[i+1:]...,
			)
			return s.Save()
		}
	}
	return errors.New("record not found")
}

func (s VaultService) Search(query string) []*Record {
	q := strings.ToLower(query)
	var res []*Record

	for _, r := range s.vault.Records {
		if strings.Contains(strings.ToLower(r.Title), q) {
			res = append(res, r)
			continue
		}
		switch r.Type {
		case RecordNote:
			if strings.Contains(strings.ToLower(r.Note.Text), q) {
				res = append(res, r)
			}
		case RecordCredential:
			if strings.Contains(strings.ToLower(r.Credential.Website), q) {
				res = append(res, r)

			}
		}
	}
	return res
}

func (s *VaultService) Export() ([]byte, error) {
	return s.storage.Read("vault.dat")
}

func (s *VaultService) Import(data []byte) error {
	_, err := s.crypto.Decrypt(s.key, data)
	if err != nil {
		return errors.New("invalid vault file or wrong password")
	}
	return s.storage.Write("vault.dat", data)
}
