package v1

import (
	"net/http"
	"practice7/internal/entity"
	usecase "practice7/internal/usercase"
	"practice7/pkg/logger"
	"practice7/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	t usecase.UserInterface
	l logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface, limiter *utils.RateLimiter) {
	r := &userRoutes{t, l}
	h := handler.Group("/users")
	{
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		// Protected routes — JWT middleware first, then rate limiter uses userID
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		protected.Use(utils.RateLimitMiddleware(limiter))
		{
			protected.GET("/me", r.GetMe)
			protected.PATCH("/promote/:id", utils.RoleMiddleware("admin"), r.PromoteUser)
		}
	}
}

// RegisterUser
func (r *userRoutes) RegisterUser(c *gin.Context) {
	var createUserDTO entity.CreateUserDTO
	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	hashedPassword, err := utils.HashPassword(createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to hash password"))
		return
	}

	role := "user"
	if createUserDTO.Role != "" {
		role = createUserDTO.Role
	}

	user := entity.User{
		Username: createUserDTO.Username,
		Email:    createUserDTO.Email,
		Password: hashedPassword,
		Role:     role,
		Verified: false,
	}

	createdUser, err := r.t.RegisterUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":       createdUser.ID,
			"username": createdUser.Username,
			"email":    createdUser.Email,
			"role":     createdUser.Role,
		},
	})
}

// LoginUser
func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	token, err := r.t.LoginUser(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid username or password"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Returns the information about the authenticated user
func (r *userRoutes) GetMe(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User ID not found in token"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID format"))
		return
	}

	user, err := r.t.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User information retrieved",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"verified": user.Verified,
		},
	})
}

// Changes the role of a specific user from user to admin
// Only accessible by admin users (RoleMiddleware ensures this)
func (r *userRoutes) PromoteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID format"))
		return
	}

	err = r.t.PromoteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated to admin",
		"user_id": userID,
	})
}
