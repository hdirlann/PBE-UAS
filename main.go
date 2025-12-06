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
	"github.com/joho/godotenv"
)

func main() {
	// load .env (optional)
	_ = godotenv.Load()

	// load config
	env := config.LoadEnv()

	// context for DB connections and shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect Postgres (this function must set database.PostgresDB for compatibility)
	if err := database.ConnectPostgres(env); err != nil {
		log.Fatalf("failed connect postgres: %v", err)
	}
	// ensure postgres closed on exit (will be attempted again in shutdown)
	defer func() {
		if database.PostgresDB != nil {
			_ = database.PostgresDB.Close()
		}
	}()

	// connect Mongo (expects ConnectMongo(ctx, env) which sets database.MongoClient & database.MongoDB)
	if err := database.ConnectMongo(ctx, env); err != nil {
		log.Fatalf("failed connect mongo: %v", err)
	}
	// ensure mongo disconnect on exit
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

	// health-check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// register routes (these functions should exist under package route)
	// If your register functions need DB references, adjust signatures accordingly.
	route.RegisterPsqlRoutes(app)
	route.RegisterMongoRoutes(app)

	// start server in goroutine
	addr := ":" + env.AppPort
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", addr)
		serverErr <- app.Listen(addr)
	}()

	// wait for interrupt or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case sig := <-quit:
		log.Printf("signal %v received, shutting down...", sig)
	case err := <-serverErr:
		if err != nil {
			log.Printf("server error: %v", err)
		}
	}

	// graceful shutdown
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
