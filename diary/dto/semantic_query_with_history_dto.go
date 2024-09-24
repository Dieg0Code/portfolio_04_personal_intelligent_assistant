package dto

type SemanticQueryWithHistoryDTO struct {
	Query   string         `json:"query"`
	History ChatHistoryDTO `json:"history"`
}
