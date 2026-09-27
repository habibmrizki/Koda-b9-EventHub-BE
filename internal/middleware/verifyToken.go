package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/habibmrizki/BE-EventHub/internal/models"
	"github.com/habibmrizki/BE-EventHub/internal/utils"
	"github.com/habibmrizki/BE-EventHub/pkg"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			log.Println("Authorization header is missing")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
				Msg:       "Unauthorized",
				Err:       "Silahkan login terlebih dahulu",
			})
			return
		}

		parts := strings.Fields(bearerToken)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			log.Println("Invalid Authorization header format. Expected: Bearer <token>")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
				Msg:       "Unauthorized",
				Err:       "Format authorization header tidak valid",
			})
			return
		}

		token := parts[1]
		if token == "" {
			log.Println("Token is empty after Bearer")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
				Msg:       "Unauthorized",
				Err:       "Silahkan login terlebih dahulu",
			})
			return
		}

		var claims pkg.Claims
		if err := claims.VerifyToken(token); err != nil {
			if strings.Contains(err.Error(), jwt.ErrTokenInvalidIssuer.Error()) || strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
				log.Println("JWT Error.\nCause: ", err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{

					IsSuccess: false,
					Code:      http.StatusUnauthorized,
					Msg:       "Unauthorized",

					Err: "Token telah kedaluwarsa atau tidak valid, silahkan login kembali",
				})
				return
			}

			log.Println("Internal Server Error.\nCause: ", err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusInternalServerError,
				Msg:       "Internal Server Error",
				Err:       "Terjadi kesalahan saat memverifikasi token",
			})
			return
		}

		if utils.IsBlacklisted(token) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusUnauthorized,
				Msg:       "Unauthorized",
				Err:       "Token sudah logout dan tidak dapat digunakan lagi",
			})
			return
		}

		ctx.Set("claims", &claims)
		ctx.Set("user_id", claims.UserId)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleVal, exists := ctx.Get("role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusForbidden,
				Msg:       "Forbidden",
				Err:       "Akses ditolak: role tidak ditemukan",
			})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				IsSuccess: false,
				Code:      http.StatusForbidden,
				Msg:       "Forbidden",
				Err:       "Akses ditolak: format role tidak valid",
			})
			return
		}

		for _, allowed := range allowedRoles {
			if strings.EqualFold(userRole, allowed) {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
			IsSuccess: false,
			Code:      http.StatusForbidden,
			Msg:       "Forbidden",
			Err:       "Akses ditolak: Anda tidak memiliki izin untuk mengakses resource ini",
		})
	}
}
