package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/services"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetAdminDashboard godoc
// @Summary      Get admin global statistics
// @Description  Mengambil metrik agregat global sistem (total user, total event, active events, total communities, global fill rate)
// @Tags         Admin Dashboard
// @Accept       json
// @Produce      json
// @Security     JWTtoken
// @Success      200  {object}  models.ResponseData
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /admin/dashboard [get]
func (h *DashboardHandler) GetAdminDashboard(ctx *gin.Context) {
	result, err := h.service.GetAdminDashboard(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal mengambil dashboard admin",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true, Code: http.StatusOK, Msg: "Dashboard admin berhasil diambil",
		Data: result,
	})
}
