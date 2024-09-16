package model

import "time"

// DiarySearchResult represents the result of a semantic search in the diary
type DiarySearchResult struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	Similarity float64   `json:"similarity"`
}
