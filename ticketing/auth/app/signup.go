package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticketing/auth/database"
)

type signUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}

func (a *app) signUp(c *gin.Context) {
	var request signUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, &ErrValidation{err})
		return
	}
	if exists, _ := a.db.IsUserExists(request.Email); exists {
		abortWithError(c, ErrBadRequest{errors.New("Email in use")})
		return
	}
	user, _ := a.db.CreateUser(database.User{
		Email:    request.Email,
		Password: request.Password,
	})
	c.JSON(http.StatusCreated, user)
	return

}
