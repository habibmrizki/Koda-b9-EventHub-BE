package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
)

type CommunityHandler struct {
	communityService *services.CommunityService
}

func NewCommunityHandler(communityService *services.CommunityService) *CommunityHandler {
	return &CommunityHandler{communityService: communityService}
}

// GetAllCommunities godoc
// @Summary      Get list of communities
// @Description  Mengambil daftar komunitas dengan fitur pencarian dan filter kategori/status
// @Tags         Communities
// @Accept       json
// @Produce      json
// @Param        search    query    string  false  "Keyword pencarian nama/deskripsi komunitas"
// @Param        category  query    string  false  "Kategori komunitas (e.g. Technology, Design, All)"
// @Param        status    query    string  false  "Status komunitas (e.g. active, all)"
// @Param        page      query    int     false  "Nomor halaman (default: 1)"
// @Param        limit     query    int     false  "Jumlah data per halaman (default: 10)"
// @Success      200       {object}  models.ResponseData{data=dto.CommunityListData}
// @Failure      500       {object}  models.ErrorResponse
// @Router       /communities [get]
func (h *CommunityHandler) GetAllCommunities(ctx *gin.Context) {
	var filter dto.CommunityFilterQuery
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusBadRequest, Msg: "Invalid query parameters",
			Err: err.Error(),
		})
		return
	}

	currentUserID := 0
	if uid, err := utils.GetUserFromCtx(ctx); err == nil {
		currentUserID = uid
	}

	result, err := h.communityService.GetAllCommunities(ctx.Request.Context(), filter, currentUserID)
	if err != nil {
		log.Println("Error GetAllCommunities:", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal mengambil daftar komunitas",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true, Code: http.StatusOK, Msg: "Berhasil mengambil daftar komunitas",
		Data: result,
	})
}

// GetCommunityByID godoc
// @Summary      Get community detail by ID
// @Description  Mengambil informasi detail komunitas beserta member dan daftar diskusi
// @Tags         Communities
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Community ID"
// @Success      200  {object}  models.ResponseData{data=dto.CommunityDetailResponse}
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /communities/{id} [get]
func (h *CommunityHandler) GetCommunityByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	communityID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusBadRequest, Msg: "Invalid community ID",
			Err: err.Error(),
		})
		return
	}

	currentUserID := 0
	if uid, err := utils.GetUserFromCtx(ctx); err == nil {
		currentUserID = uid
	}

	community, err := h.communityService.GetCommunityByID(ctx.Request.Context(), communityID, currentUserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusNotFound, Msg: "Komunitas tidak ditemukan",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true, Code: http.StatusOK, Msg: "Detail komunitas berhasil diambil",
		Data: community,
	})
}

// JoinCommunity godoc
// @Summary      Join a community
// @Description  Bergabung menjadi member komunitas
// @Tags         Communities
// @Accept       json
// @Produce      json
// @Security     JWTtoken
// @Param        id   path      int  true  "Community ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /communities/{id}/join [post]
func (h *CommunityHandler) JoinCommunity(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusUnauthorized, Msg: "Unauthorized",
			Err: err.Error(),
		})
		return
	}

	idParam := ctx.Param("id")
	communityID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusBadRequest, Msg: "Invalid community ID",
			Err: err.Error(),
		})
		return
	}

	if err := h.communityService.JoinCommunity(ctx.Request.Context(), communityID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal bergabung dengan komunitas",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil bergabung ke komunitas",
	})
}

// LeaveCommunity godoc
// @Summary      Leave a community
// @Description  Keluar dari komunitas
// @Tags         Communities
// @Accept       json
// @Produce      json
// @Security     JWTtoken
// @Param        id   path      int  true  "Community ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /communities/{id}/leave [post]
func (h *CommunityHandler) LeaveCommunity(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusUnauthorized, Msg: "Unauthorized",
			Err: err.Error(),
		})
		return
	}

	idParam := ctx.Param("id")
	communityID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusBadRequest, Msg: "Invalid community ID",
			Err: err.Error(),
		})
		return
	}

	if err := h.communityService.LeaveCommunity(ctx.Request.Context(), communityID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal keluar dari komunitas",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil keluar dari komunitas",
	})
}

// AddDiscussion godoc
// @Summary      Add a discussion post in community
// @Description  Menambahkan postingan diskusi baru pada forum komunitas
// @Tags         Communities
// @Accept       json
// @Produce      json
// @Security     JWTtoken
// @Param        id       path      int                          true  "Community ID"
// @Param        request  body      dto.CreateDiscussionRequest  true  "Discussion payload"
// @Success      201      {object}  models.Response
// @Failure      400      {object}  models.ErrorResponse
// @Failure      401      {object}  models.ErrorResponse
// @Failure      500      {object}  models.ErrorResponse
// @Router       /communities/{id}/discussions [post]
func (h *CommunityHandler) AddDiscussion(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{IsSuccess: false, Code: http.StatusUnauthorized, Msg: "Unauthorized",
			Err: err.Error(),
		})
		return
	}

	idParam := ctx.Param("id")
	communityID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{IsSuccess: false, Code: http.StatusBadRequest, Msg: "Invalid community ID",
			Err: err.Error(),
		})
		return
	}

	var req dto.CreateDiscussionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{IsSuccess: false, Code: http.StatusBadRequest, Msg: "Content diskusi wajib diisi",
			Err: err.Error(),
		})
		return
	}

	if err := h.communityService.AddDiscussion(ctx.Request.Context(), communityID, userID, req.Content); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal menambahkan diskusi",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		IsSuccess: true,
		Code:      http.StatusCreated,
		Msg:       "Diskusi berhasil ditambahkan",
	})
}

// GetPopularCommunities godoc
// @Summary      Get popular communities
// @Description  Mengambil daftar komunitas terpopuler
// @Tags         Communities
// @Accept       json
// @Produce      json
// @Param        limit  query    int  false  "Limit data (default: 5)"
// @Success      200    {object}  models.ResponseData
// @Failure      500    {object}  models.ErrorResponse
// @Router       /communities/popular [get]
func (h *CommunityHandler) GetPopularCommunities(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 5
	}

	result, err := h.communityService.GetPopularCommunities(ctx.Request.Context(), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal mengambil daftar komunitas populer",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true, Code: http.StatusOK, Msg: "Berhasil mengambil komunitas populer",
		Data: result,
	})
}

// GetCommunityMembers godoc
// @Summary      Get community members
// @Description  Mengambil daftar anggota suatu komunitas
// @Tags         Communities
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Community ID"
// @Success      200  {object}  models.ResponseData
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /communities/{id}/members [get]
func (h *CommunityHandler) GetCommunityMembers(ctx *gin.Context) {
	idParam := ctx.Param("id")
	communityID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusBadRequest, Msg: "Invalid community ID",
			Err: err.Error(),
		})
		return
	}

	members, err := h.communityService.GetCommunityMembers(ctx.Request.Context(), communityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false, Code: http.StatusInternalServerError, Msg: "Gagal mengambil daftar anggota komunitas",
			Err: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true, Code: http.StatusOK, Msg: "Daftar anggota komunitas berhasil diambil",
		Data: members,
	})
}
