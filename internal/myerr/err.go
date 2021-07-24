package myerr

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// http err
func SendHttpErr(c *gin.Context) {
	http.NotFound(c.Writer, c.Request)
	return
}
