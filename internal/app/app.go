package app

import (
	"log"
	"pr_reviewer_service/internal/config"
	"pr_reviewer_service/internal/handler"
	"pr_reviewer_service/internal/repository"
	"pr_reviewer_service/internal/usecase"
	"pr_reviewer_service/migrations"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Run - запуск сервера
func Run() {
	server := gin.Default()

	logger, err := zap.NewDevelopment()
	if err != nil {
		return
	}

	cfg, err := config.New(logger.Named("config"))
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
		return
	}

	repo, err := repository.New(cfg, logger)
	if err != nil {
		logger.Fatal("error loading repository", zap.Error(err))
		return
	}

	allMigrations := migrations.New(cfg, logger)
	if err := allMigrations.RunMigrations(); err != nil {
		logger.Fatal("error running migrations", zap.Error(err))
		return
	}

	useCase := usecase.New(repo)

	prHandler := handler.New(useCase)

	server.POST("/team/add", prHandler.CreateTeam)
	server.GET("/team/get", prHandler.GetTeam)

	server.POST("/users/setIsActive", prHandler.SetIsActive)
	server.GET("/users/getReview", prHandler.GetReview)

	server.POST("/pullRequest/create", prHandler.PullRequestCreate)
	server.POST("/pullRequest/merge", prHandler.MergePR)
	server.POST("/pullRequest/reassign", prHandler.ReassignPrReviewer)

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
