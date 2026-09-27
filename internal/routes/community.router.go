package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/handlers"
	"github.com/habibmrizki/BE-EventHub/internal/middleware"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCommunityRouter(router *gin.Engine, db *pgxpool.Pool) {
	communityRepo := repositories.NewCommunityRepository(db)
	communityService := services.NewCommunityService(communityRepo)
	communityHandler := handlers.NewCommunityHandler(communityService)

	communityRoutes := router.Group("/communities")

	// Public routes
	communityRoutes.GET("", communityHandler.GetAllCommunities)
	communityRoutes.GET("/popular", communityHandler.GetPopularCommunities)
	communityRoutes.GET("/:id", communityHandler.GetCommunityByID)
	communityRoutes.GET("/:id/members", communityHandler.GetCommunityMembers)

	// User protected routes
	userAuthGroup := communityRoutes.Group("")
	userAuthGroup.Use(middleware.VerifyToken())
	userAuthGroup.POST("/:id/join", communityHandler.JoinCommunity)
	userAuthGroup.POST("/:id/leave", communityHandler.LeaveCommunity)
	userAuthGroup.POST("/:id/discussions", communityHandler.AddDiscussion)

	// Organizer / Admin protected routes
	organizerGroup := communityRoutes.Group("")
	organizerGroup.Use(middleware.VerifyToken(), middleware.RequireRole("organizer", "admin"))

}
