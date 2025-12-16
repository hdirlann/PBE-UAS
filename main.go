// main.go
//
// @title Alumni / Prestasi Management API
// @version 1.0
// @description API untuk mengelola prestasi / alumni (Postgres + Mongo) menggunakan Clean Architecture
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer " followed by a JWT token
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"clean-arch/config"
	"clean-arch/database"
	"clean-arch/route"

	"github.com/gofiber/fiber/v2"

	// swagger handler for fiber
	"github.com/gofiber/swagger"

	// import generated swagger docs (created by `swag init`)
	_ "clean-arch/docs"
)

func main() {
	// load environment
	env := config.LoadEnv()

	// context to control DB connect / shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect Postgres (sets database.PostgresDB)
	if err := database.ConnectPostgres(env); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	// ensure Postgres closed on exit
	defer func() {
		if database.PostgresDB != nil {
			_ = database.PostgresDB.Close()
		}
	}()

	// Connect Mongo (sets database.MongoClient & database.MongoDB)
	if err := database.ConnectMongo(ctx, env); err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	// ensure mongo disconnected on exit
	defer func() {
		if database.MongoClient != nil {
			_ = database.MongoClient.Disconnect(ctx)
		}
	}()

	// create fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("ok")
	})

	// swagger UI (after you generate docs with swag)
	// browse to: http://localhost:3000/swagger/index.html
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Register application routes (single entrypoint)
	// Replace old RegisterPsqlRoutes + RegisterMongoRoutes with RegisterAPIRoutes
	route.RegisterAPIRoutes(app)

	// start server (graceful shutdown support)
	port := env.AppPort
	if port == "" {
		port = "3000"
	}
	addr := ":" + port

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", addr)
		serverErr <- app.Listen(addr)
	}()

	// Wait for SIGINT or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case sig := <-quit:
		log.Printf("signal %v received - shutting down...", sig)
	case err := <-serverErr:
		if err != nil {
			log.Printf("server error: %v", err)
		}
	}

	// graceful shutdown with timeout
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := app.Shutdown(); err != nil {
		log.Printf("fiber shutdown error: %v", err)
	}

	// close DBs
	if database.MongoClient != nil {
		_ = database.MongoClient.Disconnect(shutdownCtx)
	}
	if database.PostgresDB != nil {
		_ = database.PostgresDB.Close()
	}

	log.Println("server stopped")
}
