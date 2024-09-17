package model

import (
	"time"

	"github.com/google/uuid"
)

type Diary struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at" time_format:"2006-01-02 15:04:05.999999999"`
	Embedding []float32 `json:"embedding"`
}
