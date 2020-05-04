package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *app) currentUser(c *gin.Context) {
	c.String(http.StatusOK, "Hi there")
}

