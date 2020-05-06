package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticketing/auth/database"
	"ticketing/auth/model"
	"ticketing/auth/services"
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

func signOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("jwt", "", -1, "/", "", false, true)
		c.String(http.StatusOK, "Logged out")
	}
}

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
