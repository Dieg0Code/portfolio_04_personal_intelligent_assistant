package model

import "github.com/google/uuid"

type UserMessage struct {
	ID             uuid.UUID `json:"id"`
	MessageContent string    `json:"message_content"`
	SenderLocation string    `json:"sender_location"`
	CreatedAt      string    `json:"created_at"`
	Embedding      []float32 `json:"embedding"`
}
