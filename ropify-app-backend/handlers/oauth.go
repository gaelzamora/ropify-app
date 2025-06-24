package handlers

import (
	"context"
	"time"

	"github.com/gaelzamora/ropify-app/config"
	"github.com/gaelzamora/ropify-app/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	service *services.OAuthService
}

// Inicia el flujo de autenticación de Google
func (h *OAuthHandler) GoogleLogin(ctx *fiber.Ctx) error {
	url := config.GoogleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return ctx.Redirect(url)
}

// Maneja la redirección de Google después de la autenticación
func (h *OAuthHandler) GoogleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No code provided",
		})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token, user, err := h.service.HandleGoogleLogin(context, code)
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
	route.Get("/google/login", handler.GoogleLogin)
	route.Get("/google/callback", handler.GoogleCallback)

	// Rutas para Facebook
	//route.Get("/facebook/login", handler.FacebookLogin)
	//route.Get("/facebook/callback", handler.FacebookCallback)

	// Rutas para Twitter
	//route.Get("/twitter/login", handler.TwitterLogin)
	//route.Get("/twitter/callback", handler.TwitterCallback)
}
