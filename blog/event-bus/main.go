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
)

const QueryService = "http://localhost:4002"

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
			c.AbortWithError(http.StatusBadGateway, err)
			return
		}

		fmt.Printf("Received: %+v\n", event)

		// Handle event
		switch event.Type {
		case "PostCreated":
			data, _ := json.Marshal(event)
			http.Post("http://localhost:4002/events", "application/json", bytes.NewReader(data))
		case "CommentCreated":
			data, _ := json.Marshal(event)
			fmt.Printf("%+v", event)
			http.Post("http://localhost:4002/events", "application/json", bytes.NewReader(data))
		default:

		}

	})

	addr := os.Getenv("APP_PORT")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
