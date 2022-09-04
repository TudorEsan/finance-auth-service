package helpers

import (
	"auth-service/models"
	"time"

	sharedmodels "github.com/TudorEsan/shared-finance-app-golang/sharedModels"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateTokens(user models.User) (string, string, error) {
	claims := &sharedmodels.SignedDetails{
		Email:    *user.Email,
		Username: *user.Email,
		Id:       user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * 15).Unix(),
		},
	}
	refreshClaims := &sharedmodels.SignedDetails{
		Id:       user.ID.Hex(),
		Email:    *user.Email,
		Username: *user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(sharedmodels.GetSecretKey())
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(sharedmodels.GetSecretKey())
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func RemoveCookies(c *gin.Context) {
	c.SetCookie("token", "", 60*60*24*30, "", "", false, false)
	c.SetCookie("refreshToken", "", 60*60*24*30, "", "", false, false)
}
