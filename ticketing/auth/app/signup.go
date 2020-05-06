package app

import (
	"errors"
	"net/http"
	"ticketing/auth/database"

	"github.com/gin-gonic/gin"
)

type signUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}

func (a *app) signUp(c *gin.Context) {
	var request signUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		abortWithError(c, &ErrBadRequest{err})
		return
	}
	if exists, _ := a.db.IsUserExists(request.Email); exists {
		abortWithError(c, ErrBadRequest{errors.New("email in use")})
		return
	}

	user, _ := a.db.CreateUser(database.User{
		Email:    request.Email,
		Password: request.Password,
	})

	jwt, _ := a.db.GetJWT(user)

	c.SetCookie("jwt", jwt, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, jwt)
	return
}
