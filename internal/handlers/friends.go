package handlers

import "github.com/gin-gonic/gin"

func RegisterFriends(r *gin.Engine) {
	r.POST("/friends", SendFriendRequest)
	r.GET("/friends", ListFriends)
}

func SendFriendRequest(c *gin.Context) {

}

func ListFriends(c *gin.Context) {

}
