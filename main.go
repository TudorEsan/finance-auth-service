package main

import (
	"auth-service/database"
	"auth-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

func main() {
	// dependencies
	r := gin.Default()
	gr := r.Group("")
	l := hclog.New(&hclog.LoggerOptions{
		Name: "AUTH",
	})
	client := database.DbInstace()

	// routes
	routes.AuthRoutes(gr, l, client)

	port := ":4001"
	r.Run(port)
}
