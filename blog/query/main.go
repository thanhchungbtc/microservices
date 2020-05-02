package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

type Post struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Comments []Comment `json:"comments"`
}
type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	PostID  string `json:"post_id"`
	Status  string `json:"status"`
}

var Posts []*Post

func init() {
	Posts = make([]*Post, 0)
}

type Event struct {
	Type string
	Data interface{}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/posts", func(c *gin.Context) {
		c.JSON(http.StatusOK, Posts)
	})

	r.POST("/events", func(c *gin.Context) {
		var event Event
		if err := c.ShouldBind(&event); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		fmt.Printf("Process event: %v\n", event.Type)

		switch event.Type {
		case "PostCreated":
			data, _ := json.Marshal(event.Data)
			var post Post
			if err := json.Unmarshal(data, &post); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			post.Comments = make([]Comment, 0)
			Posts = append(Posts, &post)
		case "CommentCreated":
			data, _ := json.Marshal(event.Data)
			var comment Comment
			if err := json.Unmarshal(data, &comment); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			// find the post
			idx := -1
			for idx = 0; idx < len(Posts); idx++ {
				if Posts[idx].ID == comment.PostID {
					break
				}
			}
			Posts[idx].Comments = append(Posts[idx].Comments, comment)
		case "CommentUpdated":
			var comment Comment
			data, _ := json.Marshal(event.Data)
			json.Unmarshal(data, &comment)

			// find the post
			idx := -1
			for idx = 0; idx < len(Posts); idx++ {
				if Posts[idx].ID == comment.PostID {
					break
				}
			}

			for i, c := range Posts[idx].Comments {
				if c.ID == comment.ID {
					Posts[idx].Comments[i] = comment
					break
				}
			}
		}
	})

	addr := os.Getenv("APP_PORT")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
