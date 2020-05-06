package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func signOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("jwt", "", -1, "/", "", false, true)
		c.String(http.StatusOK, "Logged out")
	}
}
