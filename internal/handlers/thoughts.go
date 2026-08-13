package handlers

import "github.com/gin-gonic/gin"

func RegisterThoughts(r *gin.Engine) {
	r.POST("/thoughts", SendThoughts)
	r.GET("/thoughts", ListThoughts)
}

func SendThoughts(c *gin.Context) {

}

func ListThoughts(c *gin.Context) {

}
