package main

import (
	"fmt"
	"os"

	"stayt/internal/auth"
	"stayt/internal/cloud"
	"stayt/internal/db"
	"stayt/internal/handlers"
	"stayt/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	pool, _, err := db.DbConnect()
	if err != nil {
		fmt.Println("Fauled to establish connection with the database: %w", err)
		return
	}
	defer pool.Close()
	fmt.Println("Connection with database sucessfully established!")

	bucketName := os.Getenv("AWS_S3_BUCKET")
	// Create instance of s3client
	s3Client, err := cloud.Config(bucketName)
	if err != nil {
		fmt.Println("Failed to establish connection with AWS cloud: %w", err)
		return
	}

	// Registers the router with all routes
	routes.RegisterRoutes(r, pool, s3Client)

	protected := r.Group("/")
	protected.Use(auth.Middleware())
	protected.GET("/thoughts", handlers.ListAllReceivedThoughts(pool, s3Client))
	protected.GET("/thoughts/unread", handlers.ListUnreadThoughts(pool, s3Client))
	protected.POST("/thoughts", handlers.SendThoughts(pool, s3Client))
	protected.PATCH("/thoughts/:id/viewed", handlers.UpdateThoughtViewedStatus(pool))
	protected.POST("/friends", handlers.SendFriendRequest(pool))
	protected.GET("/friends", handlers.ListFriends(pool))
	protected.PATCH("/friends/:id/accept", handlers.AcceptFriendRequest(pool))
	protected.PATCH("/friends/:id/reject", handlers.RejectFriendRequest(pool))
	protected.GET("/me", handlers.ViewUserProfile(pool))

	r.Run()
}
