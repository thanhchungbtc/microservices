package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticketing/auth/database"
)

type app struct {
	r  *gin.Engine
	db *database.Database
}

func New(db *database.Database) *app {
	r := gin.Default()
	a := &app{r: r, db: db}

	r.GET("/ping", a.ping)
	r.POST("/api/users/signin", a.signIn)
	r.POST("/api/users/signout", a.signOut)
	r.POST("/api/users/signup", a.signUp)
	r.GET("/api/users/currentUser", a.currentUser)

	return a
}

func (a *app) Run(addr string) error {
	return a.r.Run(addr)
}

func (a *app) ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func abortWithError(c *gin.Context, err error) {
	type Response struct {
		Errors []errorResponse `json:"errors"`
	}
	var response Response
	var code int

	switch err := err.(type) {
	case Error:
		code = err.StatusCode()
		response.Errors = err.Json()

	default:
		code = http.StatusInternalServerError
		response.Errors = []errorResponse{newErrorResponse("oops. Something went wrong.")}
	}

	c.AbortWithStatusJSON(code, response)

	return
}
