package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type App struct {
	r *gin.Engine
}

func New() *App {
	r := gin.Default()
	a := &App{r: r}

	r.GET("/ping", a.ping)
	r.POST("/api/users/signin", a.signIn)
	r.POST("/api/users/signout", a.signOut)
	r.POST("/api/users/signup", a.signUp)
	r.GET("/api/users/currentUser", a.currentUser)

	return a
}

func (a *App) Run(addr string) error {
	return a.r.Run(addr)
}

func (a *App) ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func abortWithError(c *gin.Context, code int, err error) {
	type errorResponse struct {
		Message string `json:"message"`
		Param   string `json:"param,omitempty"`
	}

	var errResponses []errorResponse
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errResponses = append(errResponses, errorResponse{
				Message: fmt.Sprintf("error validation for field %s with tag %s", e.Field(), e.Tag()),
				Param:   e.Field(),
			})
		}
	} else {
		errResponses = append(errResponses, errorResponse{
			Message: err.Error(),
			Param:   "",
		})
	}

	c.AbortWithStatusJSON(code, errResponses)
	return
}

func abortInternalError(c *gin.Context, err error) {
	log.Printf("Error: +%v", err)
	abortWithError(c, http.StatusInternalServerError, errors.New("oops, something went wrong"))
}
