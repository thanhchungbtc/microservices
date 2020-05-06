package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"ticketing/auth/database"
)

type app struct {
	r  *gin.Engine
	db *database.Database
}

func New(db *database.Database) *app {
	r := gin.Default()
	a := &app{r: r, db: db}

	userRepo := database.NewUserRepository(db)

	r.GET("/ping", a.ping)
	r.POST("/api/users/signin", signIn(userRepo))
	r.POST("/api/users/signout", signOut())
	r.POST("/api/users/signup", signUp(userRepo))

	r.Use(authRequired()).
		GET("/api/users/currentUser", a.currentUser)

	return a
}

func (a *app) Run(addr string) error {
	return a.r.Run(addr)
}

func (a *app) ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("jwt")
		if err != nil {
			abortWithError(c, err)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user", claims)
			c.Next()
		} else {
			abortWithError(c, err)
			return
		}
	}
}

func abortWithError(c *gin.Context, err error) {

	type Response struct {
		Errors []errorResponse `json:"errors"`
	}
	var response Response
	var code int

	switch err := err.(type) {
	case ErrBadRequest:
		code = http.StatusBadRequest
		response.Errors = err.Responses()
	default:
		switch err {
		case ErrDatabaseConnection:
			code = http.StatusInternalServerError
		default:
			code = http.StatusBadRequest
			fmt.Printf("Error: %+v\n ", err)
			debug.PrintStack()
		}
		response.Errors = []errorResponse{errorResponse{Message: err.Error()}}
	}

	c.AbortWithStatusJSON(code, response)

	return
}
