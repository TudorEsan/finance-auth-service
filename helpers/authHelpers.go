package helpers

import (
	"App/database"
	"App/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil
	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(c, userType)
	return err
}

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("unauthorized to access this resorce")
	}
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ValidUsername(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	count, err := userCollection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}
	return nil
}

func CheckPassword(dbUser models.User, user models.User) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(*dbUser.Password), []byte(*user.Password))
	if err != nil {
		return errors.New("credentials are not good")
	}
	return nil
}

func UpdateTokens(c *gin.Context, token string, refreshToken string, userId string) (models.User, error) {
	c.SetCookie("token", token, 60*60*24*30, "", "", false, false)
	c.SetCookie("refreshToken", refreshToken, 60*60*24*30, "", "", false, false)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	var user models.User
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return user, errors.New("not a valid object id")
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	fmt.Println("UserId: ", userId)
	err = userCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.D{
		{"$set", bson.D{{"refreshToken", refreshToken}}},
	}, &opts).Decode(&user)
	return user, err
}

// calculate fibonaci sequence
