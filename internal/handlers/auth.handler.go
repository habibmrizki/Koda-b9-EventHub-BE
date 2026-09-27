package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/habibmrizki/BE-EventHub/pkg"
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

// Login godoc
// @Summary      Login user
// @Description  Autentikasi pengguna dan mengembalikan JWT Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginRequest  true  "Kredensial Login"
// @Success      200      {object}  models.ResponseData
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var body dto.LoginRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Data login tidak valid",
			Err:       err.Error(),
		})
		return
	}

	authData, err := h.service.Login(ctx.Request.Context(), body)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusUnauthorized,
			Msg:       "Login gagal",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Login berhasil",
		Data:      authData,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Stateless logout
// @Tags         Auth
// @Security     JWTtoken
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	parts := strings.Fields(bearerToken)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusUnauthorized,
			Msg:       "Format authorization header tidak valid",
			Err:       "Unauthorized",
		})
		return
	}

	token := parts[1]
	var claims *pkg.Claims
	if claimsAny, exists := ctx.Get("claims"); exists {
		if c, ok := claimsAny.(*pkg.Claims); ok {
			claims = c
		}
	}

	if err := h.service.Logout(ctx.Request.Context(), token, claims); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal logout",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Logout berhasil",
	})
}

// ForgotPassword godoc
// @Summary      Forgot Password request
// @Description  Meminta verifikasi email untuk reset password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ForgotPasswordRequest  true  "Email yang terdaftar"
// @Success      200      {object}  models.Response
// @Failure      400      {object}  models.ErrorResponse
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var body dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Format email tidak valid",
			Err:       err.Error(),
		})
		return
	}

	if err := h.service.ForgotPassword(ctx.Request.Context(), body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Gagal memproses lupa password",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Permintaan reset password berhasil diproses",
	})
}

// ResetPassword godoc
// @Summary      Reset / Create New Password
// @Description  Mengubah password user setelah verifikasi lupa password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ResetPasswordRequest  true  "Data Password Baru"
// @Success      200      {object}  models.Response
// @Failure      400      {object}  models.ErrorResponse
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(ctx *gin.Context) {
	var body dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Data reset password tidak valid",
			Err:       err.Error(),
		})
		return
	}

	if err := h.service.ResetPassword(ctx.Request.Context(), body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Gagal memperbarui password",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Password berhasil diperbarui, silakan login",
	})
}
