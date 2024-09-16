package model

type Diary struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	Embedding []float32 `json:"embedding"`
}
