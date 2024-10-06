package service

import "github.com/dieg0code/rag-diary/diary/dto"

type DiaryService interface {
	CreateDiary(diary dto.CreateDiaryDTO) error
	SaveUserMessage(userMessage string, ip string) error
	// GetDiary(id int) (*dto.DiaryDTO, error)
	// GetAllDiaries() ([]*dto.DiaryDTO, error)
	// DeleteDiary(id int) error
	SematicSearch(query string) (string, error)
	RAGResponse(query dto.SemanticQueryWithHistoryDTO) (string, error)
}
