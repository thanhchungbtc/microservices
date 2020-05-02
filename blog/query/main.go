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
	ID       string `json:"id"`
	Title    string `json:"title"`
	Comments []Comment `json:"comments"`
}
type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	PostID  string `json:"post_id"`
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
		switch event.Type {
		case "PostCreated":
			data, _ := json.Marshal(event.Data)
			fmt.Println(data)
			var post Post
			if err := json.Unmarshal(data, &post); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			post.Comments = make([]Comment, 0)
			Posts = append(Posts, &post)
		case "CommentCreated":
			data, _ := json.Marshal(event.Data)
			fmt.Printf("%v", event.Data)
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
			fmt.Println(idx)
			Posts[idx].Comments = append(Posts[idx].Comments, comment)
		}
	})

	addr := os.Getenv("APP_PORT")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
