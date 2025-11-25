package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"clean-arch/app/database"
	"clean-arch/app/route"
)

func main() {
	// load .env (optional)
	_ = godotenv.Load()

	// connect postgres first
	if err := database.ConnectPostgres(); err != nil {
		log.Fatal("postgres connect:", err)
	}
	// connect mongo
	if err := database.ConnectMongo(); err != nil {
		log.Fatal("mongo connect:", err)
	}

	app := fiber.New()
	// logger middleware
	app.Use(logger.New())

	// register routes (psql & mongo)
	route.RegisterPsqlRoutes(app)
	route.RegisterMongoRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
