package main

import (
	"fmt"
	"log"
	"os"
	"ticketing/auth/app"
)

func main() {
	a := app.New()

	addr := os.Getenv("APP_PORT")
	if addr == "" {
		addr = "4000"
	}
	fmt.Println("Start listening at port " + addr)
	if err := a.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
