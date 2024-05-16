package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func health() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "UP")
		return
	}
}
