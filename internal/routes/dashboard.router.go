package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/handlers"
	"github.com/habibmrizki/BE-EventHub/internal/middleware"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDashboardRouter(router *gin.Engine, db *pgxpool.Pool) {
	repo := repositories.NewDashboardRepository(db)
	service := services.NewDashboardService(repo)
	handler := handlers.NewDashboardHandler(service)

	// Admin Dashboard
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.VerifyToken(), middleware.RequireRole("admin"))
	adminGroup.GET("/dashboard", handler.GetAdminDashboard)
}
