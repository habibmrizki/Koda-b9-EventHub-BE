package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// GetMyNotifications godoc
// @Summary      Get user notifications
// @Description  Mengambil daftar notifikasi milik user yang sedang login
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     JWTtoken
// @Success      200  {object}  models.ResponseData
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /notifications [get]
func (h *NotificationHandler) GetMyNotifications(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusUnauthorized, Msg: "Unauthorized",
			Err: err.Error(),
		})
		return
	}

	notifs, err := h.service.GetMyNotifications(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal mengambil daftar notifikasi",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true, Code: http.StatusOK, Msg: "Notifikasi berhasil diambil",
		Data: notifs,
	})
}
