package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/handlers"
	"github.com/habibmrizki/BE-EventHub/internal/middleware"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitNotificationRouter(router *gin.Engine, db *pgxpool.Pool) {
	repo := repositories.NewNotificationRepository(db)
	service := services.NewNotificationService(repo)
	handler := handlers.NewNotificationHandler(service)

	notifGroup := router.Group("/notifications")
	notifGroup.Use(middleware.VerifyToken())

	notifGroup.GET("", handler.GetMyNotifications)
}
