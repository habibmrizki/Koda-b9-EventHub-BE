package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/handlers"
	"github.com/habibmrizki/BE-EventHub/internal/middleware"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitEventRouter(router *gin.Engine, db *pgxpool.Pool) {
	eventRepository := repositories.NewEventRepository(db)
	eventHandler := handlers.NewEventHandler(eventRepository)

	// Public Event Routes
	eventGroup := router.Group("/events")
	eventGroup.GET("", eventHandler.GetEvents)
	eventGroup.GET("/upcoming", eventHandler.GetUpcomingEvents)
	eventGroup.GET("/:id", eventHandler.GetEventDetail)
	eventGroup.POST("/:id/join", middleware.VerifyToken(), eventHandler.JoinEvent)
	eventGroup.POST("/:id/leave", middleware.VerifyToken(), eventHandler.LeaveEvent)

	// Organizer routes
	organizerGroup := eventGroup.Group("")
	organizerGroup.Use(middleware.VerifyToken(), middleware.RequireRole("organizer"))
	organizerGroup.POST("", eventHandler.CreateEvent)
	organizerGroup.PUT("/:id", eventHandler.UpdateEvent)

	// User Event Routes (My Events)
	userGroup := router.Group("/user")
	userGroup.Use(middleware.VerifyToken())
	userGroup.GET("/events", eventHandler.GetMyEvents)

}
