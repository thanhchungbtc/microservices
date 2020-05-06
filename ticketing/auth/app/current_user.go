package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *app) currentUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"currentUser": user,
	})
}
