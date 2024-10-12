package data

import (
	"errors"

	"github.com/dieg0code/rag-diary/diary/model"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

type DiaryRepositoryImpl struct {
	supabase *supabase.Client
}

// insertUserMessage implements DiaryRepository.
func (d *DiaryRepositoryImpl) InsertUserMessage(UserMessage *model.UserMessage) error {
	_, count, err := d.supabase.From("user_message").Insert(UserMessage, false, "", "representation", "exact").Execute()
	if err != nil {
		logrus.WithError(err).Error("cannot insert user message")
		return err
	}

	if count == 0 {
		logrus.Warn("user message not inserted")
		return errors.New("user message not inserted")
	}

	return nil
}

// SemanticSearch implements DiaryRepository.
func (d *DiaryRepositoryImpl) SemanticSearch(queryEmbedding []float32, similarityThreshold float32, matchCount int) (string, error) {
	params := map[string]interface{}{
		"query_embedding":      queryEmbedding,
		"similarity_threshold": similarityThreshold,
		"match_count":          matchCount,
	}

	response := d.supabase.Rpc("search_diary", "exact", params)
	logrus.WithField("response", response).Info("semantic search response")

	if response == "" {
		logrus.Error("cannot get semantic search")
		return "", errors.New("cannot get semantic search")
	}

	return response, nil
}

// InsertDiary implements DiaryRepository.
func (d *DiaryRepositoryImpl) InsertDiary(diary *model.Diary) error {
	_, count, err := d.supabase.From("diary").Insert(diary, false, "", "representation", "exact").Execute()
	if err != nil {
		logrus.WithError(err).Error("cannot insert diary")
		return err
	}

	if count == 0 {
		logrus.Warn("diary not inserted")
		return errors.New("diary not inserted")
	}

	return nil
}

func NewDiaryRepositoryImpl(supabase *supabase.Client) DiaryRepository {
	return &DiaryRepositoryImpl{
		supabase: supabase,
	}
}
