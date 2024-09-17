package data

import (
	"encoding/json"
	"errors"

	"github.com/dieg0code/rag-diary/diary/model"
	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

type DiaryRepositoryImpl struct {
	supabase *supabase.Client
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

// DeleteDiary implements DiaryRepository.
// func (d *DiaryRepositoryImpl) DeleteDiary(id int) error {
// 	idStr := strconv.Itoa(id)
// 	_, count, err := d.supabase.From("diary").Delete("*", "exact").Eq("id", idStr).Execute()
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot delete diary")
// 		return err
// 	}

// 	if count == 0 {
// 		logrus.WithField("id", id).Warn("diary not found")
// 		return errors.New("diary not found")
// 	}
// 	return nil
// }

// GetAllDiaries implements DiaryRepository.
// func (d *DiaryRepositoryImpl) GetAllDiaries() ([]*model.Diary, error) {
// 	data, count, err := d.supabase.From("diary").Select("*", "exact", false).Execute()
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot get all diaries")
// 		return nil, errors.New("error getting registries from database")
// 	}

// 	if count == 0 {
// 		logrus.Warn("no diaries found")
// 		return nil, nil
// 	}

// 	var diaries []*model.Diary
// 	err = json.Unmarshal(data, &diaries)
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot unmarshal diaries")
// 		return nil, err
// 	}

// 	return diaries, nil
// }

// GetDiary implements DiaryRepository.
// func (d *DiaryRepositoryImpl) GetDiary(id int) (*model.Diary, error) {
// 	data, count, err := d.supabase.From("diary").Select("*", "exact", false).Eq("id", strconv.Itoa(id)).Execute()
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot get diary")
// 		return nil, err
// 	}

// 	if count == 0 {
// 		logrus.WithField("id", id).Warn("diary not found")
// 		return nil, errors.New("diary not found")
// 	}

// 	var diaries []*model.Diary
// 	err = json.Unmarshal(data, &diaries)
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot unmarshal diary")
// 		return nil, err
// 	}

// 	return diaries[0], nil
// }

// InsertDiary implements DiaryRepository.
func (d *DiaryRepositoryImpl) InsertDiary(diary *model.Diary) (*model.Diary, error) {
	data, count, err := d.supabase.From("diary").Insert(diary, false, "", "representation", "exact").Execute()
	if err != nil {
		logrus.WithError(err).Error("cannot insert diary")
		return nil, err
	}

	if count == 0 {
		logrus.Warn("diary not inserted")
		return nil, errors.New("diary not inserted")
	}

	var diaries []*model.Diary
	err = json.Unmarshal(data, &diaries)
	if err != nil {
		logrus.WithError(err).Error("cannot unmarshal diary")
		return nil, err
	}

	return diaries[0], nil
}

func NewDiaryRepositoryImpl(supabase *supabase.Client) DiaryRepository {
	return &DiaryRepositoryImpl{
		supabase: supabase,
	}
}
