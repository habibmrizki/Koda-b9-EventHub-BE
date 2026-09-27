package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/handlers"
	"github.com/habibmrizki/BE-EventHub/internal/middleware"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitTestimonialRouter(router *gin.Engine, db *pgxpool.Pool) {
	repo := repositories.NewTestimonialRepository(db)
	service := services.NewTestimonialService(repo)
	handler := handlers.NewTestimonialHandler(service)

	testimonialGroup := router.Group("/testimonials")

	testimonialGroup.GET("", handler.GetTestimonials)
	testimonialGroup.POST("", middleware.VerifyToken(), handler.CreateTestimonial)
}
