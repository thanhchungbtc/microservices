package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *App) currentUser(c *gin.Context) {
	c.String(http.StatusOK, "Hi there")
}

