package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
)

type TestimonialHandler struct {
	service *services.TestimonialService
}

func NewTestimonialHandler(service *services.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{service: service}
}

// GetTestimonials godoc
// @Summary      Get testimonials
// @Description  Mengambil daftar testimoni featured/semua untuk landing page
// @Tags         Testimonials
// @Accept       json
// @Produce      json
// @Param        limit  query    int   false  "Limit data (default: 6)"
// @Param        all    query    bool  false  "Ambil semua testimonial jika true"
// @Success      200    {object}  models.ResponseData
// @Failure      500    {object}  models.ErrorResponse
// @Router       /testimonials [get]
func (h *TestimonialHandler) GetTestimonials(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "6")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 6
	}

	all := ctx.Query("all") == "true"

	var testimonials []dto.TestimonialResponse
	if all {
		testimonials, err = h.service.GetAllTestimonials(ctx.Request.Context())
	} else {
		testimonials, err = h.service.GetFeaturedTestimonials(ctx.Request.Context(), limit)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal mengambil daftar testimoni",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true, Code: http.StatusOK, Msg: "Berhasil mengambil testimoni",
		Data: testimonials,
	})
}

// CreateTestimonial godoc
// @Summary      Create testimonial
// @Description  Menambahkan testimoni baru oleh user yang sedang login
// @Tags         Testimonials
// @Accept       json
// @Produce      json
// @Security     JWTtoken
// @Param        request  body      dto.CreateTestimonialRequest  true  "Testimonial payload"
// @Success      201      {object}  models.ResponseData
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /testimonials [post]
func (h *TestimonialHandler) CreateTestimonial(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusUnauthorized, Msg: "Unauthorized",
			Err: err.Error(),
		})
		return
	}

	var req dto.CreateTestimonialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusBadRequest, Msg: "Data testimoni tidak valid",
			Err: err.Error(),
		})
		return
	}

	result, err := h.service.CreateTestimonial(ctx.Request.Context(), userID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal menyimpan testimoni",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseData{
		IsSuccess: true, Code: http.StatusCreated, Msg: "Testimoni berhasil ditambahkan",
		Data: result,
	})
}
