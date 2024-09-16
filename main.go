package main

import (
	"net/http"

	"github.com/dieg0code/rag-diary/db"
	"github.com/dieg0code/rag-diary/diary/controller"
	"github.com/dieg0code/rag-diary/diary/data"
	"github.com/dieg0code/rag-diary/diary/provider"
	"github.com/dieg0code/rag-diary/diary/service"
	"github.com/dieg0code/rag-diary/router"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Error("cannot load env file")
	}

	db := db.NewDBConnection()
	repo := data.NewDiaryRepositoryImpl(db)
	openai := provider.NewOperAiClient()
	service := service.NewDiaryServiceImpl(openai, repo)
	controller := controller.NewDiaryControllerImpl(service)

	r := router.NewRouter(controller)

	ginRouter := r.InitRoutes()

	serer := &http.Server{
		Addr:    ":8080",
		Handler: ginRouter,
	}

	err = serer.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Error("cannot start server")
	}

	logrus.Info("server started")
}
