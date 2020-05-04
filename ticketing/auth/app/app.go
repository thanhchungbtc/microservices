package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type handler struct {
	r *gin.Engine
}

func New() *handler {
	r := gin.Default()
	h := &handler{r: r}

	r.GET("/ping", h.ping)
	r.POST("/api/users/signin", h.signIn)
	r.POST("/api/users/signout", h.signOut)
	r.POST("/api/users/signup", h.signUp)
	r.GET("/api/users/currentUser", h.currentUser)

	return h
}

func (h *handler) Run(addr string) error {
	return h.r.Run(addr)
}

func (h *handler) ping(c *gin.Context) {
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
