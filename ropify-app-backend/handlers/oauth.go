package handlers

import (
	"context"
	"time"

	"github.com/gaelzamora/ropify-app/services"
	"github.com/gofiber/fiber/v2"
)

type OAuthHandler struct {
	service *services.OAuthService
}

func (h *OAuthHandler) GoogleToken(ctx *fiber.Ctx) error {
	var request struct {
		AccessToken string `json:"access_token"`
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request body",
		})
	}

	// Usar el access_token para obtener información del usuario
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, user, err := h.service.HandleGoogleToken(context, request.AccessToken)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully authenticated with Google",
		"data": fiber.Map{
			"token": token,
			"user":  user,
		},
	})
}

// Implementar métodos similares para Facebook y Twitter

func NewOAuthHandler(route fiber.Router, service *services.OAuthService) {
	handler := &OAuthHandler{
		service: service,
	}

	// Rutas para Google
	route.Post("/google/token", handler.GoogleToken)

	// Rutas para Facebook
	//route.Get("/facebook/login", handler.FacebookLogin)
	//route.Get("/facebook/callback", handler.FacebookCallback)

	// Rutas para Twitter
	//route.Get("/twitter/login", handler.TwitterLogin)
	//route.Get("/twitter/callback", handler.TwitterCallback)
}
