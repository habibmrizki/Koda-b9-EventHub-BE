package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/BE-EventHub/pkg"
)

func GetUserFromCtx(c *gin.Context) (int, error) {
	claims, ok := c.Get("claims")
	if !ok {
		return 0, errors.New("claims not found in context, token might be missing")
	}

	userClaims, ok := claims.(*pkg.Claims)
	if !ok {
		return 0, errors.New("invalid claims format")
	}

	return userClaims.UserId, nil
}
