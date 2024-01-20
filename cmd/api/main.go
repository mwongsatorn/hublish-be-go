package main

import (
	"fmt"
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5713",
		AllowCredentials: true,
	}))
	database.ConnectDatabase()
	routes.RegisterRoutes(app)
	port := os.Getenv("PORT")
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
