package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func genPostID() string {
	token := make([]byte, 4)
	rand.Read(token)
	return base64.URLEncoding.EncodeToString(token)
}

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Event struct {
	Type string
	Data interface{}
}

var Posts []*Post

func init() {
	Posts = []*Post{
		&Post{
			ID:    genPostID(),
			Title: "Hello World",
		},
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/posts", func(c *gin.Context) {
		c.JSON(http.StatusOK, Posts)
	})

	r.POST("/posts/create", func(c *gin.Context) {
		var post Post
		if err := c.ShouldBind(&post); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		post.ID = genPostID()
		Posts = append(Posts, &post)

		// emit event
		eventData := &Event{
			Type: "PostCreated",
			Data: post,
		}
		data, _ := json.Marshal(eventData)
		go http.Post("http://event-bus-srv:4005/events", "application/json", bytes.NewReader(data))

		c.JSON(http.StatusCreated, post)
	})

	r.POST("/events", func(c *gin.Context) {
		fmt.Println("Received")
	})

	addr := os.Getenv("APP_PORT")

	fmt.Println("v55")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
