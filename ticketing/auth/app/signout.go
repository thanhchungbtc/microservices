package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *app) signOut(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.String(http.StatusOK, "Logged out")
}
