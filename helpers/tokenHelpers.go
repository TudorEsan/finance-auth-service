package helpers

import (
	"auth-service/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	Email    string
	Username string
	Id       string
	jwt.StandardClaims
}

var SECRET_KEY []byte = getSecretKey()

func GenerateTokens(user models.User) (string, string, error) {
	claims := &SignedDetails{
		Email:    *user.Email,
		Username: *user.Email,
		Id:       user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * 15).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		Id:       user.ID.Hex(),
		Email:    *user.Email,
		Username: *user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func RemoveCookies(c *gin.Context) {
	c.SetCookie("token", "", 60*60*24*30, "", "", false, false)
	c.SetCookie("refreshToken", "", 60*60*24*30, "", "", false, false)
}

func getSecretKey() []byte {

	secret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		panic("JWT_SECRET not present")
	}
	return []byte(secret)
}
