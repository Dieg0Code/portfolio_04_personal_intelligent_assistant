package controller

import "github.com/gin-gonic/gin"

type DiaryController interface {
	CreateDiary(c *gin.Context)
	GetDiary(c *gin.Context)
	GetAllDiaries(c *gin.Context)
	DeleteDiary(c *gin.Context)
	SemanticSearch(c *gin.Context)
	RAGResponse(c *gin.Context)
}
