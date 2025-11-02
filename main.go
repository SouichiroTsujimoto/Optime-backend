package main

import (
	"log"
	"os"

	"yotei-backend/database"
	"yotei-backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// 開発環境では.envファイルを読み込む
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	}

	// データベースに接続
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// マイグレーションを実行
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Fiberアプリを作成
	app := fiber.New(fiber.Config{
		AppName: "Yotei Backend API v1.0.0",
	})

	// ミドルウェアを設定
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// ヘルスチェックエンドポイント
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

	// APIルート
	api := app.Group("/api/v1")

	// イベント関連のエンドポイント
	api.Post("/events", handlers.CreateEvent)
	// api.Get("/events", handlers.GetAllEvents)
	api.Get("/events/:id", handlers.GetEvent)

	// サーバーを起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
