package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(repo *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{
		service: services.NewAuthService(repo),
	}
}

// Register godoc
// @Summary      Register new user
// @Description  Mendaftarkan pengguna baru ke sistem
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "Data Registrasi"
// @Success      201      {object}  models.ResponseData
// @Failure      400      {object}  models.ErrorResponse
// @Failure      409      {object}  models.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var body dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Data validasi register tidak sesuai",
			Err:       err.Error(),
		})
		return
	}

	user, err := h.service.Register(ctx.Request.Context(), body)
	if err != nil {
		if errors.Is(err, repositories.ErrEmailExists) {
			ctx.JSON(http.StatusConflict, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusConflict,
				Msg:       "Email sudah terdaftar",
				Err:       err.Error(),
			})
			return
		}

		log.Printf("[ERROR] Failed to register user")
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal mendaftar (Terjadi kesalahan pada server)",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusCreated,
		Msg:       "Registrasi berhasil",
		Data:      user,
	})
}
