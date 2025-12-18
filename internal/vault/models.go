package vault

import "time"

type Meta struct {
	Salt     []byte `json:"salt"`
	Verifier []byte `json:"verifier"`
}

type Vault struct {
	Records []Record `json:"records"`
}

type Record struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
