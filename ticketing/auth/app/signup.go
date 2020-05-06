package app

import (
	"errors"
	"net/http"
	"ticketing/auth/database"
	"ticketing/auth/model"

	"github.com/gin-gonic/gin"
)

type signUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}

func signUp(userRepo *database.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request signUpRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			abortWithError(c, &ErrBadRequest{err})
			return
		}
		if exists, _ := userRepo.Exists(request.Email); exists {
			abortWithError(c, ErrBadRequest{errors.New("email in use")})
			return
		}

		user, _ := userRepo.Create(model.User{
			Email:    request.Email,
			Password: request.Password,
		})

		jwt, _ := userRepo.GetJWT(user)

		c.SetCookie("jwt", jwt, 3600, "/", "", false, true)
		c.JSON(http.StatusOK, jwt)
		return
	}
}
