package main

import (
	"log"

	"github.com/ayo-ajayi/teamsync/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err:=godotenv.Load(); err!=nil {
		log.Fatal("Error loading .env file ", err.Error())
	}
	app.NewApp(":8080", app.Router()).Start()
}