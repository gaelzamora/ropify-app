package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
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
	// Obtener la URI de redirección si se proporciona (para apps móviles)
	redirectUri := ctx.Query("redirect_uri", "")

	// Crear un estado con la URI de redirección codificada
	state := "default_state"
	if redirectUri != "" {
		stateData := map[string]string{"redirect_uri": redirectUri}
		stateBytes, _ := json.Marshal(stateData)
		state = base64.StdEncoding.EncodeToString(stateBytes)
	}

	url := config.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return ctx.Redirect(url)
}

// Maneja la redirección de Google después de la autenticación
func (h *OAuthHandler) GoogleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	stateParam := ctx.Query("state")

	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No code provided",
		})
	}

	// Procesar estado y extraer redirect_uri si existe
	var redirectUri string
	if stateParam != "default_state" {
		stateBytes, err := base64.StdEncoding.DecodeString(stateParam)
		if err == nil {
			var stateData map[string]string
			if json.Unmarshal(stateBytes, &stateData) == nil {
				redirectUri = stateData["redirect_uri"]
			}
		}
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

	// Si hay una URI de redirección móvil, redirigir ahí con los datos
	if redirectUri != "" {
		responseData := map[string]interface{}{
			"token":   token,
			"user_id": user.ID,
		}
		jsonData, _ := json.Marshal(responseData)
		encodedData := base64.URLEncoding.EncodeToString(jsonData)

		redirectTo := redirectUri
		if strings.Contains(redirectUri, "?") {
			redirectTo += "&data=" + encodedData
		} else {
			redirectTo += "?data=" + encodedData
		}

		return ctx.Redirect(redirectTo)
	}

	// Respuesta JSON normal para clientes web
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Successfully authenticated with Google",
		"data": fiber.Map{
			"token": token,
			"user":  user,
		},
	})
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
	route.Get("/google/login", handler.GoogleLogin)
	route.Get("/google/callback", handler.GoogleCallback)
	route.Post("/google/token", handler.GoogleToken)

	// Rutas para Facebook
	//route.Get("/facebook/login", handler.FacebookLogin)
	//route.Get("/facebook/callback", handler.FacebookCallback)

	// Rutas para Twitter
	//route.Get("/twitter/login", handler.TwitterLogin)
	//route.Get("/twitter/callback", handler.TwitterCallback)
}
