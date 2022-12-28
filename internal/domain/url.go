package domain

import "github.com/google/uuid"

type URL struct {
	ID       uuid.UUID `json:"id"`
	Original string    `json:"original"`
	Short    string    `json:"short"`
}
