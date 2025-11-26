package app

import (
	"pr_reviewer_service/internal/config"
	"pr_reviewer_service/internal/handler"
	"pr_reviewer_service/internal/middleware"
	"pr_reviewer_service/internal/repository"
	"pr_reviewer_service/internal/usecase"
	"pr_reviewer_service/migrations"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Run - запуск сервера
func Run() {
	server := gin.Default()

	server.Use(middleware.PrometheusMiddleware())

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

	//Teams
	teamGroup := server.Group("/team")
	teamGroup.POST("/add", prHandler.CreateTeam)
	teamGroup.GET("/get", prHandler.GetTeam)

	//Users
	usersGroup := server.Group("/users")
	usersGroup.POST("/setIsActive", prHandler.SetIsActive)
	usersGroup.GET("/getReview", prHandler.GetReview)

	//Pull Request
	prGroup := server.Group("/pullRequest")
	prGroup.POST("/create", prHandler.PullRequestCreate)
	prGroup.POST("/merge", prHandler.MergePR)
	prGroup.POST("/reassign", prHandler.ReassignPrReviewer)

	//Metrics
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))

	if err := server.Run(":8080"); err != nil {
		logger.Fatal("error running server", zap.Error(err))
	}
}
