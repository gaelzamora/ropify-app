package handlers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/gaelzamora/ropify-app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OutfitHandler struct {
	repository      models.OutfitRepository
	outfitGenerator *services.OutfitGeneratorService
}

// Crear outfit
func (h *OutfitHandler) CreateOutfit(ctx *fiber.Ctx) error {
	var outfit models.Outfit
	if err := ctx.BodyParser(&outfit); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newOutfit, err := h.repository.AddOutfit(context, &outfit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   newOutfit,
	})
}

// Editar outfit
func (h *OutfitHandler) UpdateOutfit(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	outfitID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid outfit ID",
		})
	}

	var updateData map[string]interface{}
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Convierte tags a pq.StringArray si existe
	if tags, ok := updateData["tags"]; ok {
		if arr, ok := tags.([]interface{}); ok {
			strArr := make(pq.StringArray, len(arr))
			for i, v := range arr {
				strArr[i] = v.(string)
			}
			updateData["tags"] = strArr
		}
	}

	// Convierte garment_ids a pq.StringArray si existe
	if garmentIDs, ok := updateData["garment_ids"]; ok {
		if arr, ok := garmentIDs.([]interface{}); ok {
			strArr := make(pq.StringArray, len(arr))
			for i, v := range arr {
				strArr[i] = v.(string)
			}
			updateData["garment_ids"] = strArr
		}
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updatedOutfit, err := h.repository.UpdateOutfit(context, outfitID, updateData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updatedOutfit,
	})
}

// Eliminar outfit
func (h *OutfitHandler) DeleteOutfit(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	outfitID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid outfit ID",
		})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.repository.DeleteOutfit(context, outfitID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Outfit deleted",
	})
}

// Archivar outfit
func (h *OutfitHandler) ArchiveOutfit(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	outfitID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid outfit ID",
		})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.repository.ArchiveOutfit(context, outfitID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Outfit archived",
	})
}

// Visualizar outfit por ID
func (h *OutfitHandler) GetOutfit(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	outfitID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid outfit ID",
		})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outfit, err := h.repository.GetOutfitByID(context, outfitID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Outfit not found",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   outfit,
	})
}

// Listar outfits de un usuario
func (h *OutfitHandler) GetOutfitsByUser(ctx *fiber.Ctx) error {
	userIDParam := ctx.Query("user_id", "")
	limitParam := ctx.Query("limit", "10")
	pageParam := ctx.Query("page", "1")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}
	limit, _ := strconv.Atoi(limitParam)
	page, _ := strconv.Atoi(pageParam)
	offset := (page - 1) * limit

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outfits, err := h.repository.GetOutfitsByUser(context, userID, limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   outfits,
	})
}

func (h *OutfitHandler) GenerateRandomOutfit(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outfit, garments, err := h.outfitGenerator.GenerateRandomOutfit(context, userID)
	if err != nil {
		// Separamos los errores en dos categorías: esenciales y no esenciales

		// 1. Errores por falta de prendas esenciales (tops, bottoms, sneakers)
		if strings.Contains(err.Error(), "no top garments found") ||
			strings.Contains(err.Error(), "no bottom garments found") ||
			strings.Contains(err.Error(), "no basic tops found") {

			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "success",
				"data": fiber.Map{
					"outfit":   nil,
					"garments": []*models.Garment{},
					"message":  "Not enough essential garments (tops, bottoms) to generate an outfit. Please add more clothing to your wardrobe.",
				},
			})
		}

		// 2. Errores por falta de prendas no esenciales (vestidos, accesorios, etc.)
		// Estos errores no deberían detener la generación de outfits
		if strings.Contains(err.Error(), "no dress garments found") ||
			strings.Contains(err.Error(), "no accessories found") {

			// En este caso, intentamos generar un outfit alternativo sin estas prendas
			// Nota: Esta lógica tendría que estar implementada en el servicio OutfitGeneratorService
			// para permitir generación de outfits con prendas opcionales faltantes

			// Por ahora devolvemos una respuesta que indica que faltan prendas opcionales
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "success",
				"data": fiber.Map{
					"outfit":   nil,
					"garments": []*models.Garment{},
					"message":  "Some optional garments are missing, but you can still add more variety to your wardrobe.",
				},
			})
		}

		// Para otros tipos de errores, mantener el comportamiento original
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// El resto del código permanece igual...
	if ctx.Query("save", "false") == "true" {
		savedOutfit, err := h.repository.AddOutfit(context, outfit)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "partial_success",
				"message": "Outfit generated but not saved",
				"data": fiber.Map{
					"outfit":   outfit,
					"garments": garments,
				},
			})
		}
		outfit = savedOutfit
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"outfit":   outfit,
			"garments": garments,
		},
	})
}

func NewOutfitHandler(router fiber.Router, repository models.OutfitRepository, garmentRepository models.GarmentRepository) {
	outfitGenerator := services.NewOutfitGeneratorService(garmentRepository)

	handler := &OutfitHandler{
		repository:      repository,
		outfitGenerator: outfitGenerator,
	}

	router.Post("/", handler.CreateOutfit)

	router.Patch("/:id", handler.UpdateOutfit)
	router.Patch("/:id/archive", handler.ArchiveOutfit)

	router.Get("/", handler.GetOutfitsByUser)
	router.Get("/generate", handler.GenerateRandomOutfit)
	router.Get("/:id", handler.GetOutfit)

	router.Delete("/:id", handler.DeleteOutfit)
}
