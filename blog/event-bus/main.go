package main

import (
	"bytes"
	"encoding/json"
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

var Events []*Event

func init() {
	Events = make([]*Event, 0)
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

		// Saving the event
		Events = append(Events, &event)

		data, _ := json.Marshal(event)
		http.Post("http://posts-clusterip-srv:4000/events", "application/json", bytes.NewReader(data))
		http.Post("http://comments-srv:4001/events", "application/json", bytes.NewReader(data))
		http.Post("http://query-srv:4002/events", "application/json", bytes.NewReader(data))
		http.Post("http://moderation-srv:4003/events", "application/json", bytes.NewReader(data))
	})

	r.GET("/events", func(c *gin.Context) {
		c.JSON(http.StatusOK, Events)
	})

	addr := os.Getenv("APP_PORT")
	if err := r.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
