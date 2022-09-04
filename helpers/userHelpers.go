package helpers

import (
	"auth-service/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func GetUserForDb(user models.User) (models.User, error) {
	// formats user to be passed to the db
	user.CreateDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	fmt.Println(user.ID)
	hashedPassw, err := HashPassword(*user.Password)
	if err != nil {
		return models.User{}, err
	}
	*user.Password = hashedPassw
	return user, nil
}

func GetUser(userCollection *mongo.Collection, id string) (user models.User, err error) {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return user, fmt.Errorf("could not find user in the db")
	}
	return
}
