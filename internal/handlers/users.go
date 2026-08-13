package handlers

import "github.com/gin-gonic/gin"

func RegisterUsers(r *gin.Engine) {
	r.GET("/me", Me)
}

func Me(c *gin.Context) {

}
