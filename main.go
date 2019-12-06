package main

import (
	"log"
	"os"

	"github.com/KevLehmann/packky-tracker-api/app"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	gin.ForceConsoleColor()
	r := gin.Default()
	app.InitRoutes(r)

	appPort := os.Getenv("APP_PORT")
	log.Println("Starting app, running server on port ", appPort)
	log.Fatal(r.Run(appPort))
}
