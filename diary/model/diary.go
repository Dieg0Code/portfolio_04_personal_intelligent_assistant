package model

import "time"

type Diary struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Embedding []float32 `json:"embedding"`
}
