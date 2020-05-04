package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) signIn(c *gin.Context) {
	c.String(http.StatusOK, "Hi there")
}