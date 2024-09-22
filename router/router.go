package router

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/dieg0code/rag-diary/diary/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	diaryController controller.DiaryController
	ginLambda       *ginadapter.GinLambda
}

func NewRouter(diaryController controller.DiaryController) *Router {
	return &Router{
		diaryController: diaryController,
	}
}

func (r *Router) InitRoutes() *Router {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://dieg0code.site"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
			// diaryRoute.GET("/:id", r.diaryController.GetDiary)
			// diaryRoute.GET("", r.diaryController.GetAllDiaries)
			// diaryRoute.DELETE("/:id", r.diaryController.DeleteDiary)
			diaryRoute.POST("/semantic-search", r.diaryController.SemanticSearch)
			diaryRoute.POST("/rag-response", r.diaryController.RAGResponse)
		}
	}

	r.ginLambda = ginadapter.New(router)
	return r
}

func (r *Router) Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return r.ginLambda.ProxyWithContext(ctx, req)
}
