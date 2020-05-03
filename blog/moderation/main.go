package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	PostID  string `json:"post_id"`
	Status  string `json:"status"`
}

type Event struct {
	Type string
	Data interface{}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/events", func(c *gin.Context) {
		var event Event
		if err := c.ShouldBind(&event); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		fmt.Printf("Received: %+v\n", event.Type)

		switch event.Type {
		case "CommentCreated":
			var comment Comment
			data, _ := json.Marshal(event.Data)
			if err := json.Unmarshal(data, &comment); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			// do moderation logic
			go func(comment Comment) {
				time.Sleep(3 * time.Second)
				if strings.Contains(comment.Content, "bad") {
					comment.Status = "reject"
				} else {
					comment.Status = "approved"
				}
				// emit event
				event := Event{
					Type: "CommentModerated",
					Data: comment,
				}
				data, _ := json.Marshal(event)
				http.Post("http://event-bus-srv:4005/events", "application/json", bytes.NewReader(data))
			}(comment)

		}
	})

	addr := os.Getenv("APP_PORT")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
