package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"stayt/internal/cloud"
	"stayt/internal/handlers/sanitizer"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type SendThoughtReq struct {
	ReceiverID  int    `json:"receiver_id" binding:"required"`
	ContentText string `json:"content_text"`
	MediaType   string `json:"media_type" binding:"required, oneof = picture doodle music voice_memo message"`
}

func RegisterThoughts(r *gin.Engine, pool *pgxpool.Pool, s3Client *cloud.S3Client) {
	r.POST("/thoughts", SendThoughts(pool, s3Client))
	r.GET("/thoughts", ListThoughts)
}

func SendThoughts(pool *pgxpool.Pool, s3Client *cloud.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// To send thought :
		// reqs:
		// 1) receiver_id must be a valid user (sql takes care of)
		// 2) curr user must have jwt token (curruserID taken from .mustget)
		// 3) media type must be 1/5 options (sql takes care of)
		// process
		// 1) parse incoming http req
		// verifies users exist at the same time
		// 2) upload thought file to different database, returning content_url
		// cloud object storage
		// 3) otherwise, store content text

		var req SendThoughtReq
		ctx := c.Request.Context()
		var contentURL *string
		currentUserID := c.MustGet("user_id").(int)

		thoughtData := c.PostForm("thoughtData")
		if thoughtData == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing thoughts metadata"})
			return
		}

		// Manually unmarshal thought data
		if err := json.Unmarshal([]byte(thoughtData), &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// If there is content text, then we have a message and no raw file
		if req.MediaType != "message" {
			// Extract separate raw file
			file, err := c.FormFile("rawFile")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			openedFile, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
				return
			}
			defer openedFile.Close()

			// Sanitize file before uploading to cloud
			sanitizedFile, err := sanitizer.Sanitize(openedFile)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
				return
			}

			bucketName := os.Getenv("AWS_S3_BUCKET")
			key := "placeholder" // TODO: write unique url generator

			// Upload to AWS cloud
			client := s3Client.Client
			_, err = client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   sanitizedFile,
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s3Client.BucketName, key)
			contentURL = &url
		}
		var contentText *string
		if req.ContentText != "" {
			contentText = &req.ContentText
		}

		// Store thought entry in database
		_, err := pool.Exec(ctx,
			"INSERT into thoughts (user_id, receiver_id, content_text, content_url, media_type) VALUES $1, $2, $3, $4, $5",
			currentUserID, req.ReceiverID, contentText, contentURL, req.MediaType)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "thought sent"})
	}
}

func ListThoughts(c *gin.Context) {

}
