package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/ostd2000/Auth/routes"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")
	routes.AuthRouter(api)
}