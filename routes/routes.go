package routes

import (
	"stayt/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	handlers.RegisterFriends(r)
	handlers.RegisterThoughts(r)
	handlers.RegisterUsers(r)
	handlers.RegisterAuthRoutes(r)
}
