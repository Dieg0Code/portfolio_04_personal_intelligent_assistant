package router

import (
	"github.com/dieg0code/rag-diary/diary/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	diaryController controller.DiaryController
}

func NewRouter(diaryController controller.DiaryController) *Router {
	return &Router{
		diaryController: diaryController,
	}
}

func (r *Router) InitRoutes() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	baseRoute := router.Group("/api/v1")
	{
		diaryRoute := baseRoute.Group("/diary")
		{
			diaryRoute.POST("", r.diaryController.CreateDiary)
			diaryRoute.GET("/:id", r.diaryController.GetDiary)
			diaryRoute.GET("", r.diaryController.GetAllDiaries)
			diaryRoute.DELETE("/:id", r.diaryController.DeleteDiary)
			diaryRoute.POST("/semantic-search", r.diaryController.SemanticSearch)
			diaryRoute.POST("/rag-response", r.diaryController.RAGResponse)
		}
	}

	return router
}
