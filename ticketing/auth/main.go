package main

import (
	"fmt"
	"log"
	"os"
	"ticketing/auth/app"
	"ticketing/auth/database"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalln("Error connect to mongoDB", err)
	}
	fmt.Println("Connect to DB successfully.")

	a := app.New(db)

	addr := os.Getenv("APP_PORT")
	if addr == "" {
		addr = "4000"
	}
	fmt.Println("Start listening at port " + addr)

	if err := a.Run(":" + addr); err != nil {
		log.Fatal("Posts failed starting")
	}
}
