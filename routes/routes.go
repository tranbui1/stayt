package routes

import (
	"stayt/internal/cloud"
	"stayt/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool, s3Client *cloud.S3Client) {
	handlers.RegisterFriends(r, pool)
	handlers.RegisterThoughts(r, pool, s3Client)
	handlers.RegisterUsers(r, pool)
	handlers.RegisterAuthRoutes(r, pool)
}
