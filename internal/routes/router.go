package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/docs"
	"github.com/habibmrizki/BE-EventHub/internal/middleware"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	// inisialisasi engine gin
	router := gin.Default()
	router.Use(middleware.CORSMiddleware)

	// Static route untuk file upload (public folder)
	router.Static("/public", "./public")
	router.Static("/uploads", "./uploads")

	// swaggo configuration
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// setup routing
	InitAuthRouter(router, db)
	InitUserRouter(router, db)
	InitEventRouter(router, db)
	InitCommunityRouter(router, db)
	InitTestimonialRouter(router, db)
	InitNotificationRouter(router, db)
	InitDashboardRouter(router, db)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.Response{
			IsSuccess: false,
			Code:      http.StatusNotFound,
			Msg:       "Page not found!",
		})
	})
	return router
}
