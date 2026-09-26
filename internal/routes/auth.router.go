package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/handlers"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("/register", authHandler.Register)

}
