package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/internal/dto"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/repositories"
	"github.com/habibmrizki/BE-EventHub/internal/services"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
)

type EventHandler struct {
	service *services.EventService
}

func NewEventHandler(repo *repositories.EventRepository) *EventHandler {
	return &EventHandler{
		service: services.NewEventService(repo),
	}
}

// GetEvents godoc
// @Summary      Get list of events with search & filter
// @Description  Mengambil daftar event yang dipublish dengan filter pencarian, kategori, lokasi, dan pagination
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        search        query     string  false  "Keyword pencarian judul/deskripsi/overview"
// @Param        category      query     string  false  "Kategori event (Tech, Workshop, dll)"
// @Param        location      query     string  false  "Kota atau alamat lokasi event"
// @Param        locationType  query     string  false  "Tipe lokasi ('online' atau 'offline')"
// @Param        page          query     int     false  "Nomor halaman (default: 1)"
// @Param        limit         query     int     false  "Jumlah item per halaman (default: 10)"
// @Success      200           {object}  models.ResponseData
// @Failure      500           {object}  models.ErrorResponse
// @Router       /events [get]
func (h *EventHandler) GetEvents(ctx *gin.Context) {
	var filter dto.EventFilterQuery
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Parameter filter tidak valid",
			Err:       err.Error(),
		})
		return
	}

	userID, _ := utils.GetUserFromCtx(ctx)

	events, total, err := h.service.GetEvents(ctx.Request.Context(), filter, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal mengambil daftar event",
			Err:       err.Error(),
		})
		return
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Page:      page,
		Msg:       "Berhasil mengambil daftar event",
		Data: dto.EventListData{
			Events: events,
			Total:  total,
			Page:   page,
			Limit:  limit,
		},
	})
}

// GetEventDetail godoc
// @Summary      Get event detail by ID
// @Description  Mendapatkan detail lengkap event berdasarkan ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.ResponseData
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /events/{id} [get]
func (h *EventHandler) GetEventDetail(ctx *gin.Context) {
	idParam := ctx.Param("id")
	eventID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "ID event tidak valid",
			Err:       "Invalid ID parameter",
		})
		return
	}

	userID, _ := utils.GetUserFromCtx(ctx)

	event, err := h.service.GetEventDetail(ctx.Request.Context(), eventID, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{

				IsSuccess: false,
				Code:      http.StatusNotFound,
				Msg:       "Event tidak ditemukan",

				Err: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal mengambil detail event",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil mengambil detail event",

		Data: event,
	})
}

// GetUpcomingEvents godoc
// @Summary      Get upcoming events
// @Description  Mengambil daftar event yang akan datang (start_time >= now)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "Limit data (default: 6)"
// @Success      200    {object}  models.ResponseData
// @Failure      500    {object}  models.ErrorResponse
// @Router       /events/upcoming [get]
func (h *EventHandler) GetUpcomingEvents(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "6")
	limit, _ := strconv.Atoi(limitStr)

	userID, _ := utils.GetUserFromCtx(ctx)

	events, err := h.service.GetUpcomingEvents(ctx.Request.Context(), limit, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal mengambil upcoming events",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil mengambil upcoming events",

		Data: events,
	})
}

// GetMyEvents godoc
// @Summary      Get user registered events
// @Description  Mengambil daftar event yang telah didaftari oleh user yang sedang login
// @Tags         Events
// @Security     JWTtoken
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ResponseData
// @Failure      401  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /user/events [get]
func (h *EventHandler) GetMyEvents(ctx *gin.Context) {
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

	events, err := h.service.GetMyEvents(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal mengambil event saya",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil mengambil event saya",
		Data:      events,
	})
}

// JoinEvent godoc
// @Summary      Join an event
// @Description  Mendaftarkan user ke dalam suatu event (disimpan ke event_registrations)
// @Tags         Events
// @Security     JWTtoken
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      409  {object}  models.ErrorResponse
// @Router       /events/{id}/join [post]
func (h *EventHandler) JoinEvent(ctx *gin.Context) {
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

	idParam := ctx.Param("id")
	eventID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "ID event tidak valid",
			Err:       "Invalid ID parameter",
		})
		return
	}

	if err := h.service.JoinEvent(ctx.Request.Context(), userID, eventID); err != nil {
		if errors.Is(err, repositories.ErrAlreadyJoined) {
			ctx.JSON(http.StatusConflict, models.ErrorResponse{

				IsSuccess: false,
				Code:      http.StatusConflict,
				Msg:       "Anda sudah bergabung pada event ini",

				Err: err.Error(),
			})
			return
		}
		if errors.Is(err, repositories.ErrEventFull) {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{

				IsSuccess: false,
				Code:      http.StatusBadRequest,
				Msg:       "Kapasitas event sudah penuh",

				Err: err.Error(),
			})
			return
		}
		if errors.Is(err, repositories.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{

				IsSuccess: false,
				Code:      http.StatusNotFound,
				Msg:       "Event tidak ditemukan",

				Err: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal bergabung ke event",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil bergabung ke event",
	})
}

// LeaveEvent godoc
// @Summary      Leave / Cancel registration for an event
// @Description  Membatalkan keikutsertaan user dari event (dihapus dari event_registrations)
// @Tags         Events
// @Security     JWTtoken
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /events/{id}/leave [post]
func (h *EventHandler) LeaveEvent(ctx *gin.Context) {
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

	idParam := ctx.Param("id")
	eventID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "ID event tidak valid",
			Err:       "Invalid ID parameter",
		})
		return
	}

	if err := h.service.LeaveEvent(ctx.Request.Context(), userID, eventID); err != nil {
		if errors.Is(err, repositories.ErrNotJoinedEvent) {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{

				IsSuccess: false,
				Code:      http.StatusBadRequest,
				Msg:       "Anda belum terdaftar pada event ini",

				Err: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal membatalkan keikutsertaan event",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Berhasil membatalkan keikutsertaan event",
	})
}

// CreateEvent godoc
// @Summary      Create a new event (Admin / Organizer)
// @Description  Membuat event baru beserta upload thumbnail menggunakan multipart/form-data
// @Tags         Events
// @Security     JWTtoken
// @Accept       multipart/form-data
// @Produce      json
// @Param        title         formData  string  true   "Judul Event"
// @Param        overview      formData  string  false  "Ringkasan Event"
// @Param        description   formData  string  true   "Deskripsi Lengkap Event"
// @Param        eventDate     formData  string  true   "Tanggal Event (YYYY-MM-DD)"
// @Param        startTime     formData  string  true   "Waktu Mulai (HH:mm)"
// @Param        endTime       formData  string  true   "Waktu Selesai (HH:mm)"
// @Param        locationType  formData  string  true   "Tipe Lokasi ('online' atau 'offline')"
// @Param        city          formData  string  false  "Kota Lokasi Event"
// @Param        address       formData  string  false  "Alamat Lengkap Event"
// @Param        capacity      formData  int     true   "Kapasitas Peserta (min: 1)"
// @Param        communityId   formData  int     false  "ID Komunitas Penyelenggara"
// @Param        thumbnail     formData  file    false  "File Banner/Thumbnail Event (maks 5MB)"
// @Success      201           {object}  models.ResponseData
// @Failure      400           {object}  models.ErrorResponse
// @Failure      401           {object}  models.ErrorResponse
// @Failure      403           {object}  models.ErrorResponse
// @Router       /events [post]
func (h *EventHandler) CreateEvent(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusUnauthorized,
			Err:       "Unauthorized",
		})
		return
	}

	var body dto.CreateEventRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Form data event tidak valid",
			Err:       err.Error(),
		})
		return
	}

	thumbnailURL := ""
	file, err := ctx.FormFile("thumbnail")
	if err == nil && file != nil {
		filename, err := utils.FileUpload(ctx, file, "event")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
				Msg:       "Gagal mengunggah thumbnail",
				Err:       err.Error(),
			})
			return
		}
		thumbnailURL = "/public/" + filename
	}

	event, err := h.service.CreateEvent(ctx.Request.Context(), userID, body, thumbnailURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusInternalServerError,
			Msg:       "Gagal membuat event",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusCreated,
		Msg:       "Event berhasil dibuat",
		Data:      event,
	})
}

// UpdateEvent godoc
// @Summary      Update an event (Admin / Organizer)
// @Description  Memperbarui informasi event secara parsial
// @Tags         Events
// @Security     JWTtoken
// @Accept       multipart/form-data
// @Accept       json
// @Produce      json
// @Param        id            path      int     true   "Event ID"
// @Param        title         formData  string  false  "Judul Event"
// @Param        overview      formData  string  false  "Ringkasan Event"
// @Param        description   formData  string  false  "Deskripsi Lengkap"
// @Param        eventDate     formData  string  false  "Tanggal Event"
// @Param        startTime     formData  string  false  "Waktu Mulai"
// @Param        endTime       formData  string  false  "Waktu Selesai"
// @Param        locationType  formData  string  false  "Tipe Lokasi"
// @Param        city          formData  string  false  "Kota Lokasi"
// @Param        address       formData  string  false  "Alamat"
// @Param        capacity      formData  int     false  "Kapasitas"
// @Param        status        formData  string  false  "Status ('published', 'draft', dll)"
// @Param        thumbnail     formData  file    false  "File Banner/Thumbnail Baru (maks 5MB)"
// @Success      200           {object}  models.ResponseData
// @Failure      400           {object}  models.ErrorResponse
// @Failure      401           {object}  models.ErrorResponse
// @Failure      403           {object}  models.ErrorResponse
// @Router       /events/{id} [patch]
func (h *EventHandler) UpdateEvent(ctx *gin.Context) {
	userID, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusUnauthorized,
			Err:       "Unauthorized",
		})
		return
	}

	idParam := ctx.Param("id")
	eventID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "ID event tidak valid",
			Err:       "Invalid ID parameter",
		})
		return
	}

	var body dto.UpdateEventRequest
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Form data event tidak valid",
			Err:       err.Error(),
		})
		return
	}

	var thumbnailURL *string
	file, err := ctx.FormFile("thumbnail")
	if err == nil && file != nil {
		filename, err := utils.FileUpload(ctx, file, "event")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusBadRequest,
				Msg:       "Gagal mengunggah thumbnail",
				Err:       err.Error(),
			})
			return
		}
		path := "/public/" + filename
		thumbnailURL = &path
	}

	event, err := h.service.UpdateEvent(ctx.Request.Context(), eventID, userID, body, thumbnailURL)
	if err != nil {
		if errors.Is(err, repositories.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusForbidden,
				Msg:       "Akses ditolak: Anda bukan penyelenggara event ini",
				Err:       err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusBadRequest,
			Msg:       "Gagal memperbarui event",
			Err:       err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseData{
		IsSuccess: true,
		Code:      http.StatusOK,
		Msg:       "Event berhasil diperbarui",
		Data:      event,
	})
}
