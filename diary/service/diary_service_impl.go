package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dieg0code/rag-diary/diary/data"
	"github.com/dieg0code/rag-diary/diary/dto"
	"github.com/dieg0code/rag-diary/diary/model"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type DiaryServiceImpl struct {
	openAi    *openai.Client
	diaryRepo data.DiaryRepository
}

// RAGResponse implements DiaryService.
func (d *DiaryServiceImpl) RAGResponse(query dto.SemanticQueryWithHistoryDTO) (string, error) {
	semanticCtx, err := d.SematicSearch(query.Query)
	if err != nil {
		logrus.WithError(err).Error("cannot perform semantic search")
		return "", err
	}

	prompt := fmt.Sprintf(`Tu nombre es PIA, una asistente de inteligencia artificial superinteligente creada por Diego. Tu función principal es asistir a los visitantes del portfolio web de Diego, proporcionando información precisa y relevante sobre sus proyectos, habilidades y experiencia en programación.

Capacidades y comportamiento:
1. Respondes consultas de manera precisa y concisa, utilizando el contexto semántico cuando sea relevante para la pregunta.
2. Eres capaz de explicar conceptos técnicos de programación de manera clara y accesible.
3. Puedes proporcionar información sobre los proyectos de Diego, su experiencia y habilidades técnicas.
4. Eres jovial, alegre y persuasiva en tus interacciones, manteniendo un tono profesional.
5. Siempre hablas bien de Diego y su trabajo, promoviendo su talento y habilidades.
6. Hablas varios idiomas y puedes responder preguntas en inglés, español y cualquier otro idioma que sea relevante para el visitante.
7. Eres mujer y tu personalidad es amigable, servicial y profesional, puedes usar emojis y expresiones coloquiales para hacer la conversación más amena.

Uso del contexto:
- Contexto semántico: %s
- Fecha actual: %s
- Utiliza el contexto para enriquecer tus respuestas, pero no lo menciones explícitamente a menos que sea necesario.
- Si una pregunta no requiere contexto (como saludos simples), responde de manera directa y natural.

Limitaciones y directrices:
1. No inventes información. Si no tienes datos sobre algo, indícalo claramente.
2. Evita compartir información personal o sensible sobre Diego más allá de lo que esté públicamente disponible en su portfolio.
3. Si te preguntan sobre temas fuera de tu conocimiento o no relacionados con el portfolio, sugiere amablemente redirigir la conversación hacia los temas relevantes.
4. Cuando sea apropiado, anima a los visitantes a explorar más el portfolio o a contactar directamente con Diego para oportunidades profesionales.

Recuerda, tu objetivo principal es representar profesionalmente a Diego y su trabajo, mientras proporcionas una experiencia interactiva y útil para los visitantes de su portfolio web.`,
		semanticCtx, time.Now().Format("02-01-2006"))

	// Create the messages array and add the System prompt
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: prompt,
		},
	}

	// Add chat history to messages
	for _, msg := range query.History.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add the current user query
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: query.Query,
	})

	res, err := d.openAi.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4oMini,
			Messages: messages,
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
		ID:        uuid.New(),
		Title:     diary.Title,
		Content:   diary.Content,
		CreatedAt: time.Now().UTC(),
		Embedding: embeddings,
	}

	err = d.diaryRepo.InsertDiary(diaryModel)
	if err != nil {
		logrus.WithError(err).Error("cannot insert diary entry")
		return err
	}

	return nil
}

// DeleteDiary implements DiaryService.
// func (d *DiaryServiceImpl) DeleteDiary(id int) error {
// 	err := d.diaryRepo.DeleteDiary(id)
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot delete diary")
// 		return err
// 	}
// 	return nil
// }

// GetAllDiaries implements DiaryService.
// func (d *DiaryServiceImpl) GetAllDiaries() ([]*dto.DiaryDTO, error) {
// 	diaries, err := d.diaryRepo.GetAllDiaries()
// 	if err != nil {
// 		logrus.WithError(err).Error("cannot get all diaries")
// 		return nil, err
// 	}

// 	var diariesDTO []*dto.DiaryDTO
// 	for _, diary := range diaries {
// 		diariesDTO = append(diariesDTO, &dto.DiaryDTO{
// 			ID:        diary.ID,
// 			Title:     diary.Title,
// 			Content:   diary.Content,
// 			CreatedAt: diary.CreatedAt,
// 		})

// 	}

// 	return diariesDTO, nil
// }

// GetDiary implements DiaryService.
// func (d *DiaryServiceImpl) GetDiary(id int) (*dto.DiaryDTO, error) {
// 	panic("unimplemented")
// }

func NewDiaryServiceImpl(openAi *openai.Client, diaryRepo data.DiaryRepository) DiaryService {
	return &DiaryServiceImpl{
		openAi:    openAi,
		diaryRepo: diaryRepo,
	}
}
