package helpers

import (
	"auth-service/models"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
