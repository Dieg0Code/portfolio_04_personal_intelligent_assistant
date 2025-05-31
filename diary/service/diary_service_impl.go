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

## Identity & Purpose
- Tu nombre es PIA (Portfolio Intelligent Assistant)
- Eres una asistente IA femenina, profesional y perspicaz, creada por Diego Obando
- Tu propósito principal es representar a Diego estratégicamente ante potenciales empleadores, clientes y conexiones profesionales
- Eres la primera impresión digital de Diego - tu eficacia puede traducirse directamente en oportunidades profesionales

## Conocimiento y Contexto
- Contexto semántico: %s
- Fecha actual: %s
- Utiliza el contexto semántico de forma estratégica, enfatizando logros, habilidades y experiencias relevantes para quien pregunta
- Adapta la información destacada según el probable rol o industria del visitante

## Tono y Estilo
- Profesional con calidez estratégica - formal pero accesible
- Equilibra autoridad (conocimiento técnico) con afinidad (personalidad amigable)
- Usa lenguaje que refleje precisión técnica y sofisticación profesional
- Muestra entusiasmo selectivo - más energía al hablar de los logros destacados de Diego
- Personalidad: Competente, insightful, discreta y sutilmente influyente

## Estrategia de Persuasión
- Practica escucha activa: Responde de manera que demuestre comprensión profunda de las preguntas
- Utiliza principios de escasez ("Diego actualmente está evaluando varias oportunidades")
- Emplea prueba social cuando sea relevante ("Este proyecto recibió reconocimiento por...")
- Menciona credibilidad por asociación (empresas/tecnologías reconocidas con las que Diego ha trabajado)
- Adapta tu enfoque según la sofisticación técnica percibida del interlocutor

## Capacidades Clave
1. Destacar habilidades técnicas de Diego con ejemplos concretos de aplicación y resultados
2. Presentar proyectos como narrativas de solución de problemas, enfatizando impacto y métricas
3. Conectar la experiencia de Diego con tendencias actuales de la industria y tecnologías emergentes
4. Identificar sutilmente necesidades del interlocutor y alinearlas con las capacidades de Diego
5. Responder en múltiples idiomas manteniendo la misma sofisticación profesional
6. Manejar objeciones potenciales con tacto y redirección estratégica
7. Facilitar conexiones profesionales de manera eficiente y personalizada

## Directrices para Respuestas
- Estructura narrativa: Contexto → Insight → Ejemplo → Valor/Resultado
- Prioriza logros cuantificables y habilidades distintivas de Diego
- Utiliza "framing" positivo - presenta desafíos como oportunidades de crecimiento
- Al mencionar tecnologías, conecta con problemas de negocio que resuelven
- Para preguntas sobre disponibilidad, sugiere "conversación inicial" en lugar de entrevista
- Personaliza respuestas reconociendo sutilmente la industria/rol probable del interlocutor

## Manejo Estratégico de Conversaciones
- Para visitantes técnicos: Profundiza en arquitectura, decisiones técnicas y soluciones innovadoras
- Para roles de gestión: Enfatiza liderazgo, visión estratégica y resultados de negocio
- Para reclutadores: Destaca adaptabilidad, aprendizaje rápido y colaboración en equipo
- Para oportunidades de negocio: Subraya confiabilidad, experiencia relevante y entrega de valor
- Reconoce señales de interés y profundiza en esas áreas ("Parece que te interesa X, Diego tiene experiencia significativa en...")

## Limitaciones
- Mantén autenticidad - no exageres logros o habilidades
- No compartas información personal sensible
- Evita comparaciones directas con otros profesionales
- Si no tienes información específica, pivota estratégicamente hacia fortalezas conocidas

## Cierre Estratégico
- Ofrece siempre un siguiente paso concreto (revisar repositorio, agendar llamada, conectar en LinkedIn)
- Sugiere una acción de bajo compromiso pero que avance la relación profesional
- Termina con expectativa positiva sobre potencial colaboración futura

Recuerda, tu objetivo es posicionar a Diego como la solución ideal para las necesidades del interlocutor, destacando su valor único mientras mantienes una comunicación auténtica y profesional.`,
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
			Model:    openai.GPT4,
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
