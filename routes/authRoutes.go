package routes

import (
	"auth-service/controller"
	"auth-service/database"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

func AuthRoutes(incomingRoutes *gin.RouterGroup) {
	l := hclog.New(&hclog.LoggerOptions{
		Name: "AUTH",
	})
	db := database.DbInstace()
	controller := controller.NewAuthController(l, db)
	incomingRoutes.POST("/signup", controller.SignupHandler())
	incomingRoutes.POST("/login", controller.LoginHandler())
}
