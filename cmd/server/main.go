package main

import (
	"stayt/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Registers the router with all routes
	routes.RegisterRoutes(r)

	r.Run()
}
