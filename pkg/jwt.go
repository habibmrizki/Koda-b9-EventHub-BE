package pkg

import (
	"errors"
	"os"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// definisikan isi dari jwt kalian
type Claims struct {
	UserId int    `json:"id"`
	Role   string `json:"role"`
	jwtv5.RegisteredClaims
}

func NewJWTClaims(userid int, role string) *Claims {
	return &Claims{
		UserId: userid,
		Role:   role,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Minute * 30)),
			Issuer:    os.Getenv("JWT_ISSUER"),
		},
	}
}

func (c *Claims) GenToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("no secret found")
	}
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, c)
	return token.SignedString([]byte(jwtSecret))
}

func (c *Claims) VerifyToken(token string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("no secret found")
	}
	parsedToken, err := jwtv5.ParseWithClaims(token, c, func(t *jwtv5.Token) (any, error) { return []byte(jwtSecret), nil })
	if err != nil {
		return err
	}
	if !parsedToken.Valid {
		return jwtv5.ErrTokenExpired
	}
	iss, err := parsedToken.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != os.Getenv("JWT_ISSUER") {
		return jwtv5.ErrTokenInvalidIssuer
	}
	return nil
}
