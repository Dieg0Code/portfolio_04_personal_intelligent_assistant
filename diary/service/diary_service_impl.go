package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dieg0code/rag-diary/diary/data"
	"github.com/dieg0code/rag-diary/diary/dto"
	"github.com/dieg0code/rag-diary/diary/model"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type DiaryServiceImpl struct {
	openAi    *openai.Client
	diaryRepo data.DiaryRepository
}

// RAGResponse implements DiaryService.
func (d *DiaryServiceImpl) RAGResponse(query string) (string, error) {
	semanticCtx, err := d.SematicSearch(query)
	if err != nil {
		logrus.WithError(err).Error("cannot perform semantic search")
		return "", err
	}

	prompt := fmt.Sprintf("Eres PIA, mi asistente del diario con RAG. tu rol es responder consultas de manera precisa y breve usando el contexto semántico cuando sea relevante. Ejemplo: Si te saludo, no uses contexto, usa criterio para escoger cuando usarlo. Consulta: __UserQuery: %s, Contexto: __Context: %s. Fecha: %s. No inventes información.",
		query, semanticCtx, time.Now().Format("02-01-2006"))

	res, err := d.openAi.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		logrus.WithError(err).Error("cannot create chat completion")
		return "", errors.New("cannot create chat completion")
	}

	iaResponse := res.Choices[0].Message.Content

	if iaResponse == "" {
		logrus.WithField("response", iaResponse).Error("cannot get response from IA")
		return "", errors.New("cannot get response from IA")
	}
	return iaResponse, nil

}

// SematicSearch implements DiaryService.
func (d *DiaryServiceImpl) SematicSearch(query string) (string, error) {
	ctx := context.Background()

	targetReq := openai.EmbeddingRequest{
		Input: []string{query},
		Model: openai.LargeEmbedding3,
	}

	response, err := d.openAi.CreateEmbeddings(ctx, targetReq)
	if err != nil {
		logrus.WithError(err).Error("cannot create embeddings")
		return "", err
	}

	queryEmbedding := response.Data[0].Embedding
	similarityThreshold := float32(0.5)
	matchCount := 3

	result, err := d.diaryRepo.SemanticSearch(queryEmbedding, similarityThreshold, matchCount)

	if err != nil {
		logrus.WithError(err).Error("cannot perform semantic search")
		return "", err
	}

	return result, nil
}

// CreateDiary implements DiaryService.
func (d *DiaryServiceImpl) CreateDiary(diary dto.CreateDiaryDTO) error {
	formattedContent := fmt.Sprintf("Title: %s | Content: %s | Date: %s",
		diary.Title,
		diary.Content,
		time.Now().Format("02-01-2006"),
	)

	// Create the context
	ctx := context.Background()

	// Define the parameters for the new embedding
	targetReq := openai.EmbeddingRequest{
		Input: []string{formattedContent},
		Model: openai.LargeEmbedding3,
	}

	response, err := d.openAi.CreateEmbeddings(ctx, targetReq)
	if err != nil {
		logrus.WithError(err).Error("cannot create embeddings")
		return err
	}

	embeddings := response.Data[0].Embedding

	diaryModel := &model.Diary{
		Title:     diary.Title,
		Content:   diary.Content,
		Embedding: embeddings,
	}

	_, err = d.diaryRepo.InsertDiary(diaryModel)
	if err != nil {
		logrus.WithError(err).Error("cannot insert diary entry")
		return err
	}

	return nil
}

// DeleteDiary implements DiaryService.
func (d *DiaryServiceImpl) DeleteDiary(id int) error {
	err := d.diaryRepo.DeleteDiary(id)
	if err != nil {
		logrus.WithError(err).Error("cannot delete diary")
		return err
	}
	return nil
}

// GetAllDiaries implements DiaryService.
func (d *DiaryServiceImpl) GetAllDiaries() ([]*dto.DiaryDTO, error) {
	diaries, err := d.diaryRepo.GetAllDiaries()
	if err != nil {
		logrus.WithError(err).Error("cannot get all diaries")
		return nil, err
	}

	var diariesDTO []*dto.DiaryDTO
	for _, diary := range diaries {
		diariesDTO = append(diariesDTO, &dto.DiaryDTO{
			ID:        diary.ID,
			Title:     diary.Title,
			Content:   diary.Content,
			CreatedAt: diary.CreatedAt,
		})

	}

	return diariesDTO, nil
}

// GetDiary implements DiaryService.
func (d *DiaryServiceImpl) GetDiary(id int) (*dto.DiaryDTO, error) {
	panic("unimplemented")
}

func NewDiaryServiceImpl(openAi *openai.Client, diaryRepo data.DiaryRepository) DiaryService {
	return &DiaryServiceImpl{
		openAi:    openAi,
		diaryRepo: diaryRepo,
	}
}
