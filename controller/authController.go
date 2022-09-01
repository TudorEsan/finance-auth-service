package controller

// func Signup
import (
	"auth-service/database"
	"auth-service/helpers"
	helper "auth-service/helpers"
	"auth-service/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// func Login

type AuthController struct {
	l              hclog.Logger
	userCollection *mongo.Collection
}

func NewAuthController(l hclog.Logger, client *mongo.Client) *AuthController {
	collection := database.OpenCollection(client, "users", "users")

	return &AuthController{
		l,
		collection,
	}

}

func (controller *AuthController) saveUser(ctx context.Context, user models.User) error {
	_, err := controller.userCollection.InsertOne(ctx, user)
	return err
}

func (controller *AuthController) SignupHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		// check if username is not present in the database
		err := helper.ValidUsername(ctx, controller.userCollection, *user.Username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// apply logic to the user, hash password, add creation date
		userForDb, err := helper.GetUserForDb(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// insert user in the db
		err = controller.saveUser(ctx, userForDb)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// generate all the auth tokens
		jwt, refreshToken, err := helper.GenerateTokens(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		helper.SetCookies(c, jwt, refreshToken)

		c.JSON(http.StatusOK, gin.H{
			"user": userForDb,
		})
	}

}

func (controller *AuthController) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		var user models.UserLoginForm
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		err := controller.userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Username does not exist"})
			return
		}

		err = helper.CheckPassword(foundUser, user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		jwt, refreshToken, err := helper.GenerateTokens(foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Could not generate tokens"})
			return
		}

		helpers.SetCookies(c, jwt, refreshToken)
		c.JSON(http.StatusOK, foundUser)
	}
}
