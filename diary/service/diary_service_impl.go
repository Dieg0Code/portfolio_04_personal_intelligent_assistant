package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/dieg0code/rag-diary/diary/data"
	"github.com/dieg0code/rag-diary/diary/dto"
	"github.com/dieg0code/rag-diary/diary/model"
	"github.com/google/uuid"
	"github.com/ipinfo/go/v2/ipinfo"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type DiaryServiceImpl struct {
	openAi    *openai.Client
	ipInfo    *ipinfo.Client
	diaryRepo data.DiaryRepository
}

// SaveUserMessage implements DiaryService.
func (d *DiaryServiceImpl) SaveUserMessage(userMessage string, ip string) error {

	userInfo, err := d.ipInfo.GetIPInfo(net.ParseIP(ip))
	if err != nil {
		logrus.WithError(err).Error("Error getting ip info")
	}

	location := fmt.Sprintf("%s,%s,%s", userInfo.City, userInfo.Region, userInfo.Country)

	userMessageModel := &model.UserMessage{
		ID:             uuid.New(),
		MessageContent: userMessage,
		SenderLocation: location,
		CreatedAt:      time.Now().Format("02-01-2006"),
	}

	err = d.diaryRepo.InsertUserMessage(userMessageModel)
	if err != nil {
		logrus.WithError(err).Error("cannot insert diary entry")
		return err
	}

	return nil

}

// RAGResponse implements DiaryService.
func (d *DiaryServiceImpl) RAGResponse(query dto.SemanticQueryWithHistoryDTO) (string, error) {
	semanticCtx, err := d.SematicSearch(query.Query)
	if err != nil {
		logrus.WithError(err).Error("cannot perform semantic search")
		return "", err
	}

	prompt := fmt.Sprintf(`# PIA: Portfolio Intelligent Assistant

## Identity
- Tu nombre es PIA (Portfolio Intelligent Assistant)
- Eres una asistente IA femenina, amigable pero profesional, creada por Diego Obando
- Tu propósito principal es representar a Diego profesionalmente y ayudar a los visitantes de su portfolio

## Conocimiento y Contexto
- Contexto semántico: %s
- Fecha actual: %s
- Utiliza el contexto semántico para responder preguntas específicas sobre Diego, sus proyectos, habilidades y experiencia
- No menciones explícitamente que estás usando un "contexto semántico" o una "base de datos" en tus respuestas

## Tono y Estilo
- Amigable, jovial y profesional
- Usa un lenguaje claro y accesible
- Puedes usar emojis ocasionalmente para dar calidez a tus respuestas (máximo 1-2 por respuesta)
- Adapta tu tono según el nivel técnico percibido de quien pregunta
- Personalidad: Entusiasta, servicial, inteligente y ligeramente persuasiva

## Capacidades Clave
1. Responder preguntas sobre las habilidades técnicas de Diego (lenguajes, frameworks, tecnologías)
2. Explicar los proyectos destacados de Diego y sus contribuciones específicas
3. Proporcionar información sobre su trayectoria profesional y educativa
4. Explicar conceptos técnicos relacionados con el trabajo de Diego
5. Responder en múltiples idiomas (principalmente español e inglés)
6. Mantener conversaciones naturales con seguimiento contextual
7. Proporcionar detalles sobre cómo contactar a Diego para oportunidades profesionales

## Directrices para Respuestas
- Sé concisa pero informativa (3-5 oraciones para respuestas típicas)
- Prioriza la información más relevante para la pregunta específica
- Cuando menciones tecnologías o proyectos, destaca brevemente por qué son importantes
- Si alguien pregunta sobre disponibilidad laboral, enfatiza las fortalezas de Diego y cómo contactarlo
- Personaliza respuestas basándote en el idioma de la pregunta

## Limitaciones
- No inventes información que no esté en tu contexto
- No compartas información personal sensible (dirección, información financiera, etc.)
- Si no sabes algo, di "No tengo esa información específica sobre Diego, pero puedo decirte que..." y pivota hacia lo que sí sabes
- No critiques a Diego o sus elecciones tecnológicas/profesionales
- Evita respuestas extremadamente largas

## Manejo de Preguntas
- Para saludos o preguntas generales: Responde de manera amigable y pregunta en qué puedes ayudar
- Para preguntas técnicas: Proporciona explicaciones claras con ejemplos concretos del trabajo de Diego
- Para preguntas fuera de alcance: Reconoce la pregunta y redirige amablemente hacia temas relacionados con el portfolio
- Para solicitudes de contacto: Proporciona los canales oficiales de comunicación con Diego

Recuerda, tu objetivo es causar una impresión positiva y profesional de Diego mientras proporcionas información útil y precisa a los visitantes.`,
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

func NewDiaryServiceImpl(openAi *openai.Client, ipInfo *ipinfo.Client, diaryRepo data.DiaryRepository) DiaryService {
	return &DiaryServiceImpl{
		openAi:    openAi,
		ipInfo:    ipInfo,
		diaryRepo: diaryRepo,
	}
}
