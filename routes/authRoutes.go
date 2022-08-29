package routes

import (
	"go.mongodb.org/mongo-driver/"
)

func AuthRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/signup", controller.Signup())
	incomingRoutes.POST("/login", controller.Login())
}
