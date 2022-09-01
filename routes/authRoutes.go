package routes

import (
	"auth-service/controller"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthRoutes(incomingRoutes *gin.RouterGroup, l hclog.Logger, client *mongo.Client) {
	controller := controller.NewAuthController(l, client)
	incomingRoutes.POST("/signup", controller.SignupHandler())
	incomingRoutes.POST("/login", controller.LoginHandler())
}
