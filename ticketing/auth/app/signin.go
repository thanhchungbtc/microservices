package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticketing/auth/database"
	"ticketing/auth/services"
)

func (a *app) signIn(c *gin.Context) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		abortWithError(c, err)
		return
	}

	var user *database.User
	var err error
	if user, err = a.db.GetUser(req.Email); err != nil {
		abortWithError(c, ErrBadRequest{errors.New("invalid user")})
		return
	}

	passwordService := services.Password{}
	if err := passwordService.Compare(user.Password, req.Password); err != nil {
		abortWithError(c, err)
		return
	}

	jwt, _ := a.db.GetJWT(user)

	c.SetCookie("jwt", jwt, 3600, "/", "", false, true)
	c.JSON(http.StatusCreated, jwt)
	return
}
