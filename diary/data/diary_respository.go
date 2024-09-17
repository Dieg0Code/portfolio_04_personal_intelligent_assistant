package data

import "github.com/dieg0code/rag-diary/diary/model"

type DiaryRepository interface {
	InsertDiary(diary *model.Diary) error
	// GetDiary(id int) (*model.Diary, error)
	// GetAllDiaries() ([]*model.Diary, error)
	// DeleteDiary(id int) error
	SemanticSearch(queryEmbedding []float32, similarityThreshold float32, matchCount int) (string, error)
}
