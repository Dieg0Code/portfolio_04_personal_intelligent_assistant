package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dieg0code/rag-diary/db"
	"github.com/dieg0code/rag-diary/diary/controller"
	"github.com/dieg0code/rag-diary/diary/data"
	"github.com/dieg0code/rag-diary/diary/provider"
	"github.com/dieg0code/rag-diary/diary/service"
	"github.com/dieg0code/rag-diary/router"
	"github.com/sirupsen/logrus"
)

var r *router.Router

func init() {
	logrus.Info("initializing router")

	// Inicializar conexión a la base de datos
	dbClient, err := db.NewDBConnection()
	if err != nil {
		logrus.Fatalf("Failed to initialize database connection: %v", err)
	}
	repo := data.NewDiaryRepositoryImpl(dbClient)
	openai := provider.NewOperAiClient()
	service := service.NewDiaryServiceImpl(openai, repo)
	controller := controller.NewDiaryControllerImpl(service)

	r = router.NewRouter(controller)
	r.InitRoutes()

	logrus.Info("Successfully initialized all components")
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logrus.Info("handling request", req.RequestContext.RequestID)
	response, err := r.Handler(ctx, req)
	if err != nil {
		logrus.Error("error handling request", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{ "message": "Internal Server Error" }`,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "https://dieg0code.site",
				"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
				"Access-Control-Allow-Headers": "Content-Type, Authorization, X-Amz-Date, X-Api-Key, X-Amz-Security-Token",
			},
		}, err
	}

	logrus.Info("Request handled successfully")
	// Ensure CORS headers are set for all responses
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}
	response.Headers["Access-Control-Allow-Origin"] = "https://dieg0code.site"
	response.Headers["Access-Control-Allow-Methods"] = "POST, GET, OPTIONS, PUT, DELETE"
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization, X-Amz-Date, X-Api-Key, X-Amz-Security-Token"

	return response, nil
}

func main() {

	logrus.Info("Starting server")
	lambda.Start(Handler)
}
