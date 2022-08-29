package helpers

import (
	"App/models"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUser(userId string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	var user models.User
	defer cancel()
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return models.User{}, err
	}
	err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserFromContext(c *gin.Context) (user models.User, err error) {
	userAny, exists := c.Get("user")
	if !exists {
		err = errors.New("key does not exist in context")
		return
	}
	user, ok := userAny.(models.User)
	if !ok {
		err = errors.New("could not convert to user")
		return
	}
	return
}
