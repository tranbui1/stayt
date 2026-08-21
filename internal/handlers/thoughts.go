package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"stayt/internal/cloud"
	"stayt/internal/handlers/sanitizer"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type SendThoughtReq struct {
	ReceiverID  int64  `json:"receiver_id" binding:"required"`
	ContentText string `json:"content_text"`
	MediaType   string `json:"media_type" binding:"required,oneof=picture doodle music 'voice memo' message"`
}

type NewThought struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	ContentText *string   `json:"content_text,omitempty"`
	ContentURL  *string   `json:"content_url,omitempty"`
	MediaType   string    `json:"media_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateThoughtReq struct {
	ID int64 `json:"id" binding:"required"`
}

func RegisterThoughts(r *gin.Engine, pool *pgxpool.Pool, s3Client *cloud.S3Client) {
	r.POST("/thoughts", SendThoughts(pool, s3Client))
	r.GET("/thoughts/unread", ListUnreadThoughts(pool, s3Client))
	r.GET("/thoughts", ListAllReceivedThoughts(pool, s3Client))
	r.PATCH("/thoughts/:id/viewed", UpdateThoughtViewedStatus(pool))
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
		var contentKey *string
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
				c.JSON(http.StatusBadRequest, gin.H{"error": "raw file required does not exist"})
				return
			}

			openedFile, err := file.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
				return
			}
			defer openedFile.Close()

			// Sanitize file before uploading to cloud
			sanitizedFile, err := sanitizer.SanitizeFile(openedFile)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
				return
			}

			bucketName := os.Getenv("AWS_S3_BUCKET")
			key := generateUniqueKey(currentUserID, file.Filename)

			// Upload to AWS cloud
			client := s3Client.Client
			_, err = client.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   bytes.NewReader(sanitizedFile),
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			contentKey = &key
		}

		var contentText *string
		if req.ContentText != "" {
			contentText = &req.ContentText
		}

		// Store thought entry in database
		_, err := pool.Exec(ctx,
			"INSERT INTO thoughts (user_id, receiver_id, content_text, content_key, media_type) VALUES ($1, $2, $3, $4, $5)",
			currentUserID, req.ReceiverID, contentText, contentKey, req.MediaType)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "thought sent"})
	}
}

func ListUnreadThoughts(pool *pgxpool.Pool, s3Client *cloud.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// List every thought that have been sent to the current user & are unread:
		// 1) Query all thoughts that have viewed status as FALSE
		// 2) Return JSON payload containing list of
		ctx := c.Request.Context()
		currUserID := c.MustGet("user_id").(int)

		rows, err := pool.Query(ctx,
			"SELECT id, user_id, content_text, content_key, media_type, created_at FROM thoughts WHERE receiver_id = $1 AND viewed = false ORDER BY created_at DESC",
			currUserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		defer rows.Close()

		var thoughts []NewThought

		// Scan each row and append NewThoughts struct to array
		for rows.Next() {
			var thought NewThought
			var contentKey *string

			if err := rows.Scan(&thought.ID, &thought.UserID, &thought.ContentText, &contentKey, &thought.MediaType, &thought.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
				return
			}

			// Generate presigned URL if necessary
			if contentKey != nil {
				presignURL, err := generatePresignedURL(ctx, s3Client, *contentKey)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
					return
				}
				thought.ContentURL = presignURL
			}
			thoughts = append(thoughts, thought)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"thoughts": thoughts})
	}
}

func UpdateThoughtViewedStatus(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Updating viewed status:
		// 1) Query thought from JSON payload into req
		// 2) Update viewed, return success

		var req UpdateThoughtReq
		ctx := c.Request.Context()
		currUserID := c.MustGet("user_id").(int)

		// Parse incoming HTTP req into Go struct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		thoughtID := req.ID

		// Only the receiver of a thought is allowed to mark it as viewed
		result, err := pool.Exec(ctx,
			"UPDATE thoughts SET viewed = true WHERE id = $1 AND receiver_id = $2",
			thoughtID, currUserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "thought not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "thought viewed status sucessfully updated"})
	}
}

func ListAllReceivedThoughts(pool *pgxpool.Pool, s3Client *cloud.S3Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// List every thought that has been sent to the current user, read or not:
		// 1) Query all thoughts where receiver_id matches current user
		// 2) Return JSON payload containing list of thoughts
		ctx := c.Request.Context()
		currUserID := c.MustGet("user_id").(int)

		rows, err := pool.Query(ctx,
			"SELECT id, user_id, content_text, content_key, media_type, created_at FROM thoughts WHERE receiver_id = $1 ORDER BY created_at DESC",
			currUserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		defer rows.Close()

		var thoughts []NewThought

		for rows.Next() {
			var thought NewThought
			var contentKey *string

			if err := rows.Scan(&thought.ID, &thought.UserID, &thought.ContentText, &contentKey, &thought.MediaType, &thought.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
				return
			}

			if contentKey != nil {
				presignURL, err := generatePresignedURL(ctx, s3Client, *contentKey)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
					return
				}
				thought.ContentURL = presignURL
			}
			thoughts = append(thoughts, thought)
		}
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"thoughts": thoughts})
	}
}
