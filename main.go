package main

import (
	"auth-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	gr := r.Group("")
	routes.AuthRoutes(gr)

	port := ":4001"
	r.Run(port)
}
