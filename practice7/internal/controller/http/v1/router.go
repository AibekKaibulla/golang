package v1

import (
	"practice7/internal/usercase"
	"practice7/pkg/logger"
	"practice7/utils"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, userUC usecase.UserInterface, l logger.Interface, limiter *utils.RateLimiter) {
	h := handler.Group("/v1")

	// Apply rate limiter to all routes
	h.Use(utils.RateLimitMiddleware(limiter))

	newUserRoutes(h, userUC, l, limiter)
}