package helpers

import (
	"App/models"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			ExpiresAt: time.Now().Local().Add(time.Minute * 60 * 24 * 30).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		Id: user.ID.Hex(),
		Email:    *user.Email,
		Username: *user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * 30).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SECRET_KEY)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil && strings.Contains(err.Error(), "expired") {
		return nil, errors.New("token expired")
	}
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func ValidateRefreshToken(refreshToken string) (models.User, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return models.User{}, errors.New("refresh token is not valid")
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return models.User{}, errors.New("token not valid")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return models.User{}, errors.New("refresh token expired")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	var user models.User
	id, err := primitive.ObjectIDFromHex(claims.Id)
	if err != nil {
		return models.User{}, errors.New("not valid object id")
	}
	err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func RemoveCookies(c *gin.Context) {
	c.SetCookie("token", "", 60*60*24*30, "", "", false, false)
	c.SetCookie("refreshToken", "", 60*60*24*30, "", "", false, false)
}

func getSecretKey() []byte {
	envs, _ := godotenv.Read(".env")
	return []byte(envs["JWT_SECRET"])
}
