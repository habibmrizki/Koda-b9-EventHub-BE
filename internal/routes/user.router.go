package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/handlers"
	"github.com/habibmrizki/BE-EventHub/internal/middleware"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitUserRouter(router *gin.Engine, db *pgxpool.Pool) {
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	userGroup := router.Group("/user")
	userGroup.Use(middleware.VerifyToken())

	userGroup.GET("/profile", userHandler.GetProfile)
	userGroup.PATCH("/profile", userHandler.UpdateProfile)
	userGroup.PATCH("/change-password", userHandler.ChangePassword)
}
