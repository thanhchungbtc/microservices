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
		comments = append(comments, &comment)
		Comments[postID] = comments

		fmt.Printf("%+v", comment)

		// edit event
		eventData := Event{
			Type: "CommentCreated",
			Data: comment,
		}
		data, _ := json.Marshal(eventData)
		http.Post("http://localhost:4005/events", "application/json", bytes.NewReader(data))

		c.JSON(http.StatusCreated, comment)
	})

	addr := os.Getenv("APP_PORT")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
