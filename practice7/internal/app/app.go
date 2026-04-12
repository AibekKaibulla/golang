package app

import (
	"practice7/internal/controller/http/v1"
	"practice7/internal/usercase"
	"practice7/internal/usercase/repo"
	"practice7/pkg/logger"
	"practice7/pkg/postgres"
	"practice7/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() error {
	db, err := postgres.New()
	if err != nil {
		return err
	}

	log := logger.New()

	userRepo := repo.NewUserRepo(db)
	userUC := usecase.NewUserUseCase(userRepo)

	gin.SetMode(gin.ReleaseMode)
	handler := gin.Default()

	limiter := utils.NewRateLimiter(5, time.Minute)

	v1.NewRouter(handler, userUC, log, limiter)

	return handler.Run(":8090")
}