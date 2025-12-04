package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	// internal packages (sesuaikan jika path berbeda)
	"clean-arch/database"
	route "clean-arch/route"
	postgreRoutes "clean-arch/route"
)

func main() {
	// load .env jika ada (tidak fatal kalau file tidak ada)
	_ = godotenv.Load()

	// context untuk startup/shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// inisialisasi database (fungsi New harus meng-set koneksi Postgres & Mongo)
	db, err := database.New(ctx)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	// jika database.New mengembalikan object dengan Close, panggil defer
	if db != nil {
		defer func() {
			_ = db.Close(context.Background())
		}()
	}

	// buat Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("ok")
	})

	// register routes
	// - Mongo routes are in package route (app/route/mongo_routes.go)
	route.RegisterMongoRoutes(app)

	// - Postgres (psql) routes are in package postgre under app/route/postgre
	postgreRoutes.RegisterPsqlRoutes(app)

	// start server (graceful shutdown)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// run server in goroutine so we can catch signals
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
		log.Printf("received signal %v, shutting down...", sig)
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

	// try to close DB (if not already closed)
	if db != nil {
		if err := db.Close(shutdownCtx); err != nil {
			log.Printf("database close error: %v", err)
		}
	}

	log.Println("server stopped")
}
