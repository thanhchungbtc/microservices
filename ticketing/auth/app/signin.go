package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticketing/auth/database"
	"ticketing/auth/model"
	"ticketing/auth/services"
)

func signIn(userRepo *database.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		type request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			abortWithError(c, err)
			return
		}

		var user *model.User
		var err error
		if user, err = userRepo.Get(req.Email); err != nil {
			abortWithError(c, ErrBadRequest{errors.New("invalid user")})
			return
		}

		passwordService := services.Password{}
		if err := passwordService.Compare(user.Password, req.Password); err != nil {
			abortWithError(c, err)
			return
		}

		jwt, _ := userRepo.GetJWT(user)

		c.SetCookie("jwt", jwt, 3600, "/", "", false, true)
		c.JSON(http.StatusCreated, jwt)
		return
	}
}
