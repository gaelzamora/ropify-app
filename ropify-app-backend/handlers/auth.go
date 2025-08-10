package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type AuthHandler struct {
	authRepository models.AuthRepository
	service        models.AuthService
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	creds := &models.AuthCredentials{}

	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	if err := ctx.BodyParser(&creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if err := validate.Struct(creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	token, user, err := h.service.Login(context, creds)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "Successfully logged in",
		"data": &fiber.Map{
			"token": token,
			"user":  user,
		},
	})
}

func (h *AuthHandler) Register(ctx *fiber.Ctx) error {
	creds := &models.AuthCredentials{}

	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	if err := ctx.BodyParser(&creds); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid request format",
		})
	}

	missing := []string{}

	if creds.Email == "" {
		missing = append(missing, "email")
	}
	if creds.Username == "" {
		missing = append(missing, "username")
	}
	if creds.Password == "" {
		missing = append(missing, "password")
	}

	if len(missing) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Missing required fields: " + strings.Join(missing, ", "),
		})
	}

	// Validar formato de email
	if !models.IsValidEmail(creds.Email) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid email format",
		})
	}

	// Verificar unicidad de email
	if _, err := h.authRepository.GetUser(ctx.Context(), "email = ?", creds.Email); err == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Email is already in use",
		})
	}

	// Verificar unicidad de username
	if _, err := h.authRepository.GetUser(ctx.Context(), "username = ?", creds.Username); err == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Username is already in use",
		})
	}
	token, user, err := h.service.Register(context, creds)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "Successfully registered",
		"data": &fiber.Map{
			"token": token,
			"user":  user,
		},
	})
}

func NewAuthHandler(route fiber.Router, service models.AuthService, authRepository models.AuthRepository) {
	handler := &AuthHandler{
		authRepository: authRepository,
		service:        service,
	}

	route.Post("/login", handler.Login)
	route.Post("/register", handler.Register)
}
