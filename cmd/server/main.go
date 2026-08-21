package main

import (
	"fmt"
	"os"

	"stayt/internal/cloud"
	"stayt/internal/db"
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

	r.Run()
}
