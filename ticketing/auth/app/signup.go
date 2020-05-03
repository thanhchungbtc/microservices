package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type signUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}

func (a *App) signUp(c *gin.Context) {
	var request signUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, request)
}
