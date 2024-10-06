package dto

type UserMessageDTO struct {
	MessageContent string    `json:"message_content"`
	SenderLocation string    `json:"sender_location"`
	Embedding      []float32 `json:"embedding"`
}
