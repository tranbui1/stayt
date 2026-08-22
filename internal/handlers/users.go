package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userInfo struct {
	Username         string    `json:"username" binding:"required"`
	CreatedAt        time.Time `json:"created_at" binding:"required"`
	ThoughtsReceived int64     `json:"thoughts_received" binding:"required"`
	ThoughtsSent     int64     `json:"thoughts_sent" binding:"required"`
}

func RegisterUsers(r *gin.Engine, pool *pgxpool.Pool) {
	r.GET("/me", ViewUserProfile(pool))
}

func ViewUserProfile(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Return account details of current user:
		// 1) Query users database by user_id from middleware
		// 2) Return account details + other stats
		// - # of thoughts sent
		// - # of thoughts received
		ctx := c.Request.Context()
		currUserID := c.MustGet("user_id")

		var userInfo userInfo
		err := pool.QueryRow(ctx,
			"SELECT username, created_at FROM users WHERE id = $1",
			currUserID).Scan(&userInfo.Username, &userInfo.CreatedAt)

		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		// Fetch the total # of thoughts sent and received
		err = pool.QueryRow(ctx,
			"SELECT COUNT(*) FROM thoughts WHERE user_id = $1",
			currUserID).Scan(&userInfo.ThoughtsReceived) // Aggregate function returns at least one row

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		err = pool.QueryRow(ctx,
			"SELECT COUNT(*) FROM thoughts WHERE receiver_id = $1",
			currUserID).Scan(&userInfo.ThoughtsSent)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_profile": userInfo})
	}
}
