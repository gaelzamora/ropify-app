package main

import (
	"fmt"

	"github.com/gaelzamora/ropify-app/config"
	"github.com/gaelzamora/ropify-app/db"
	"github.com/gaelzamora/ropify-app/handlers"
	"github.com/gaelzamora/ropify-app/middlewares"
	"github.com/gaelzamora/ropify-app/repositories"
	"github.com/gaelzamora/ropify-app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	envConfig := config.NewEnvConfig()
	db := db.Init(envConfig, db.DBMigrator)

	app := fiber.New(fiber.Config{
		AppName:      "Ropify-App",
		ServerHeader: "Fiber",
		BodyLimit:    20 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8081,http://192.168.1.68:8081",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	// Repository
	garmentRepository := repositories.NewGarmentRepository(db)
	outfitRepository := repositories.NewOutfitRepository(db)
	authRepository := repositories.NewAuthRepository(db)

	// Service
	authService := services.NewAuthService(authRepository)
	oauthService := services.NewOAuthService(authRepository)

	server := app.Group("/api")

	// Auth handler's
	handlers.NewAuthHandler(server.Group("/auth"), authService)
	handlers.NewOAuthHandler(server.Group("/oauth"), oauthService)

	// Private route to verify if user is authenticated
	privateRoutes := server.Use(middlewares.AuthProtected(db))

	handlers.NewGarmentHandler(privateRoutes.Group("/garment"), garmentRepository)
	handlers.NewOutfitHandler(privateRoutes.Group("/outfit"), outfitRepository)

	app.Listen(fmt.Sprintf("0.0.0.0:" + envConfig.ServerPort))
}
