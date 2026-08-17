package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// STRUCT DECLARATION

// Used to parse incoming http reqs into Go language
type FriendShipReq struct {
	UserID      int    `json:"user_id" binding:"required"`
	FriendID    int    `json:"friend_id" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=pending"`
	RequestedBy int    `json:"requested_by" binding:"required"`
}

func RegisterFriends(r *gin.Engine, pool *pgxpool.Pool) {
	r.POST("/friends", SendFriendRequest(pool))
	r.GET("/friends", ListFriends(pool))
	r.PATCH("/friends/:id/accept", AcceptFriendRequest(pool))
	r.PATCH("/friends/:id/reject", RejectFriendRequest(pool))
}

func SendFriendRequest(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sending a friend request :
		// 1) verifying that the 2 user ids involved exist (taken care of by sql)
		// user id also can't be the same
		// 2) sending a post request to submit the status of the friend request
		// making sure that a request can't be sent 2x
		// validated by sql, returns error if condition is broken
		var req FriendShipReq
		ctx := c.Request.Context()

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Store smaller ID in user_id
		userID := req.UserID
		friendID := req.FriendID

		if friendID < userID {
			userID, friendID = friendID, userID
		}

		_, err := pool.Exec(ctx,
			"INSERT into friendships (user_id, friend_id, requested_by, status) VALUES ($1, $2, $3, $4)",
			userID, friendID, req.RequestedBy, req.Status)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "friend request already sent"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "friend request has been sent"})
	}
}

// Shared helper for accepting/rejecting friend request
func updateFriendShipStatus(ctx context.Context, pool *pgxpool.Pool, friendShipID int, currentUserID int, newStatus string) (int64, error) {
	result, err := pool.Exec(ctx,
		"UPDATE friendships SET status = $1 WHERE id = $2 AND status = 'pending' AND requested_by != $3",
		newStatus, friendShipID, currentUserID,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func AcceptFriendRequest(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Accepting friend request :
		// 1) retrieve friend req from database based on id of friendship
		// return not found error if not
		// 2) send update req to update status to "accepted"
		ctx := c.Request.Context()

		friendShipID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friendship id"})
			return
		}

		// TODO: write middleware to retrieve currentUserID
		currentUserID := c.MustGet("user_id").(int)

		rowsAffected, err := updateFriendShipStatus(ctx, pool, friendShipID, currentUserID, "accepted")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "friend request not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
	}
}

func RejectFriendRequest(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Rejecting friend request : same as accepting but value should be updated to rejected isntead
		ctx := c.Request.Context()

		friendShipID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friendship id"})
			return
		}

		// currentUserID comes from the authenticated user's JWT
		currentUserID := c.MustGet("user_id").(int)

		rowsAffected, err := updateFriendShipStatus(ctx, pool, friendShipID, currentUserID, "rejected")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "friend request not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "friend request rejected"})
	}
}

func ListFriends(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the friend list of current user id
		// Query all rows where either the current user could be sender or receiver
		// Return friends as user struct
		ctx := c.Request.Context()

		currentUserID := c.MustGet("user_id").(int)

		rows, err := pool.Query(ctx,
			"SELECT user_id, friend_id FROM friendships WHERE user_id = $1 OR friend_id = $1 AND status = 'accepted'",
			currentUserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}
		defer rows.Close()

		// Collect all friend IDs
		var friendIDs []int

		for rows.Next() {
			var friendID int
			var userID int

			if err := rows.Scan(&userID, &friendID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
				return
			}

			if friendID == currentUserID {
				friendIDs = append(friendIDs, userID)
			} else {
				friendIDs = append(friendIDs, friendID)
			}
		}

		// Catch error after loop over stream finishes
		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"friend_ids": friendIDs})
	}
}
