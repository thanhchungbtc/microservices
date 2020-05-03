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

type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	PostID  string `json:"post_id"`
	Status  string `json:"status"`
}

var Comments map[string][]*Comment

type Event struct {
	Type string
	Data interface{}
}

func genCommentID() string {
	token := make([]byte, 4)
	rand.Read(token)
	return base64.URLEncoding.EncodeToString(token)
}

func init() {
	Comments = make(map[string][]*Comment)
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/posts/:id/comments", func(c *gin.Context) {
		postID := c.Param("id")
		var comments []*Comment
		var ok bool
		if comments, ok = Comments[postID]; !ok {
			comments = make([]*Comment, 1)
		}
		var comment Comment
		if err := c.ShouldBind(&comment); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// create new comment
		comment.ID = genCommentID()
		comment.PostID = postID
		comment.Status = "pending"
		comments = append(comments, &comment)
		Comments[postID] = comments

		// emit event
		eventData := Event{
			Type: "CommentCreated",
			Data: comment,
		}
		data, _ := json.Marshal(eventData)
		go http.Post("http://event-bus-srv:4005/events", "application/json", bytes.NewReader(data))

		c.JSON(http.StatusCreated, comment)
	})

	r.POST("/events", func(c *gin.Context) {
		var event Event
		if err := c.ShouldBind(&event); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		fmt.Printf("Process event: %v\n", event.Type)

		switch event.Type {
		case "CommentModerated":
			// do more comment related stuff here
			event := Event{
				Type: "CommentUpdated",
				Data: event.Data,
			}
			data, _ := json.Marshal(event)
			go http.Post("http://event-bus-srv:4005/events", "application/json", bytes.NewReader(data))
		}
	})
	addr := os.Getenv("APP_PORT")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
