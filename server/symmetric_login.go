package server

import (
	"log/slog"
	"net/http"
	"test-jwt-auth/constants"
	"test-jwt-auth/entities"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func (s *basicServer) symmetricGetRf(c *gin.Context) {

}

func (s *basicServer) symmetricLogin(c *gin.Context) {
	userClaims := entities.UserClaims{
		Name:  "Thong",
		Id:    "12321bei21ghh4jkh12jkrh",
		Email: "thong@email.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jwt-symmetric",
			Subject:   "12321bei21ghh4jkh12jkrh",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour)},
			NotBefore: &jwt.NumericDate{Time: time.Now()},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        "12321bei21ghh4jkh12jkrh",
		},
	}

	otps := []jwt.TokenOption{}

	jwtTk := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims, otps...)

	accessToken, err := jwtTk.SignedString([]byte(constants.HSKey))
	if err != nil {
		slog.Error("SignedString error", slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24)},
		ID:        userClaims.ID,
	}

	jwtRf := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims, otps...)

	rf, err := jwtRf.SignedString([]byte(constants.HSKey))
	if err != nil {
		slog.Error("SignedString rf error", slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": rf,
	})
}

func (s *basicServer) validateSymmetricToken(c *gin.Context) {
	headerToken := c.GetHeader("Token")
	if headerToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	tokenData := &entities.UserClaims{}

	token, err := jwt.ParseWithClaims(headerToken, tokenData, func(t *jwt.Token) (interface{}, error) {
		return []byte(constants.HSKey), nil
	})
	if err != nil {
		slog.Error("ParseWithClaims error", slog.Any("err", err), slog.String("headerToken", headerToken))
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	} else if userClaims, ok := token.Claims.(*entities.UserClaims); ok {
		c.JSON(http.StatusOK, userClaims)
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}
