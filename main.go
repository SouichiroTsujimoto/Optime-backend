package main

import (
	"log"
	"os"
	"time"

	"yotei-backend/database"
	"yotei-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	c := cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc("@every 1m", func() {
		log.Println("Running scheduled job: Checking deadlines...")
		err := handlers.CheckDeadlinesAndFinalize()
		if err != nil {
			log.Println("Failed to check and finalize deadlines:", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	log.Println("Finalize deadlines scheduler started...")
	defer c.Stop()

	app := fiber.New(fiber.Config{
		AppName: "Yotei Backend API v1.0.0",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Yotei Backend API is running",
			"status":  "ok",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	api := app.Group("/api/v1")

	api.Post("/events", handlers.CreateEvent)
	api.Get("/events/:id", handlers.GetEvent)
	api.Post("/events/:id/participant", handlers.RegisterParticipant)
	api.Put("/events/:id/settings", handlers.UpdateEventSettings)
	api.Get("/rss/:id/feed", handlers.EventRSS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
