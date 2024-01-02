package main

import (
	"fmt"
	"hublish-be-go/internal/database"
	"hublish-be-go/internal/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()
	database.ConnectDatabase()
	routes.RegisterRoutes(app)
	port := os.Getenv("PORT")
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
