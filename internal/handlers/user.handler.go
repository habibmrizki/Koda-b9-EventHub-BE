package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{
		service: services.NewUserService(repo),
	}
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Mengambil profil user yang sedang login
// @Tags         User
// @Security     JWTtoken
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ResponseData
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /user/profile [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusUnauthorized,
			Msg:       "Unauthorized",
			Err:       err.Error(),
		})
		return
	}

	profile, err := h.service.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusNotFound,
			Msg:       "User tidak ditemukan",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil mengambil profil user",
		Data:      profile,
	})
}

// UpdateProfile godoc
// @Summary      Update user profile (Partial / Multipart Form)
// @Description  Memperbarui profil pengguna secara parsial dan upload file avatar menggunakan utils.FileUpload
// @Tags         User
// @Security     JWTtoken
// @Accept       multipart/form-data
// @Accept       json
// @Produce      json
// @Param        fullName  formData  string  false  "Nama Lengkap"
// @Param        location  formData  string  false  "Lokasi Kota"
// @Param        bio       formData  string  false  "Bio User"
// @Param        avatar    formData  file    false  "File Foto Avatar (JPG/PNG/WebP, maks 2MB)"
// @Success      200       {object}  models.ResponseData
// @Failure      400       {object}  models.ErrorResponse
// @Failure      401       {object}  models.ErrorResponse
// @Router       /user/profile [patch]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusUnauthorized,
			Msg:       "Unauthorized",
			Err:       err.Error(),
		})
		return
	}

	var body dto.UpdateProfileRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Format request tidak valid",
			Err:       err.Error(),
		})
		return
	}

	// 1. Cek apakah ada file avatar yang diupload via form-data
	file, err := ctx.FormFile("avatar")
	if err == nil && file != nil {
		filename, err := utils.FileUpload(ctx, file, "avatar")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{

				IsSuccess: false,
				Code:      http.StatusBadRequest,
				Msg:       "Gagal mengunggah foto avatar",

				Err: err.Error(),
			})
			return
		}
		avatarPath := "/public/" + filename
		body.AvatarURL = &avatarPath
	}

	updatedUser, err := h.service.UpdateProfile(ctx.Request.Context(), userID, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Gagal memperbarui profil",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Profil berhasil diperbarui",
		Data:      updatedUser,
	})
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Mengubah password user dengan verifikasi password lama
// @Tags         User
// @Security     JWTtoken
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ChangePasswordRequest  true  "Data Ganti Password"
// @Success      200      {object}  models.Response
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Router       /user/change-password [patch]
func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusUnauthorized,
			Msg:       "Unauthorized",
			Err:       err.Error(),
		})
		return
	}

	var body dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Data change password tidak valid",
			Err:       err.Error(),
		})
		return
	}

	if err := h.service.ChangePassword(ctx.Request.Context(), userID, body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Gagal mengubah password",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Password berhasil diperbarui",
	})
}
