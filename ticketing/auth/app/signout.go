package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *app) signOut(c *gin.Context) {
	c.String(http.StatusOK, "Hi there")
}
