package model

import (
	"time"

	"github.com/google/uuid"
)

type Diary struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Embedding []float32 `json:"embedding"`
}
