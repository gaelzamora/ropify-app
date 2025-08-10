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
	repository        models.OutfitRepository
	garmentRepository models.GarmentRepository
	closetRepository  models.ClosetRepository
	outfitGenerator   *services.OutfitGeneratorService
}

// Crear outfit
func (h *OutfitHandler) CreateOutfit(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	var payload struct {
		ClosetID string                       `json:"closet_id"`
		Name     string                       `json:"name"`
		Garments models.GarmentOptimizedArray `json:"garments"`
	}

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request data",
		})
	}

	closetID, err := uuid.Parse(payload.ClosetID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	closet, err := h.closetRepository.GetClosetByID(context, closetID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Closet not found",
		})
	}

	if closet.UserID != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You don't have access to this closet",
		})
	}

	if payload.Name == "" {
		payload.Name = "Generate Outfit"
	}

	outfit := models.Outfit{
		ID:        uuid.New(),
		Name:      payload.Name,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	savedOutfit, err := h.repository.CreateOutfitWithGarments(context, &outfit, payload.Garments)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error creating outfit: " + err.Error(),
		})
	}

	if err := h.closetRepository.AddOutfitToCloset(context, closetID, savedOutfit.ID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "partial_success",
			"message": "Outfit created but not added to closet. Try again later",
			"data": fiber.Map{
				"outfit": savedOutfit,
			},
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Outfit created successfully",
		"data": fiber.Map{
			"outfit": savedOutfit,
		},
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

	closetID, err := uuid.Parse(ctx.Params("closetId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	// Verificar que el closet existe y pertenece al usuario
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	closet, err := h.closetRepository.GetClosetByID(context, closetID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Closet not found",
		})
	}

	if closet.UserID != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You don't have access to this closet",
		})
	}

	// Generar outfit desde el closet
	outfit, garments, err := h.outfitGenerator.GenerateRandomOutfitFromCloset(context, closetID)
	if err != nil {
		// Manejar errores específicos
		if strings.Contains(err.Error(), "no top garments found") ||
			strings.Contains(err.Error(), "no bottom garments found") ||
			strings.Contains(err.Error(), "no sneakers found") {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "success",
				"data": fiber.Map{
					"outfit":   nil,
					"garments": []*models.Garment{},
					"message":  "This closet doesn't have all required garments (top, bottom, and sneakers).",
				},
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Guardar el outfit si se solicita
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

		// Añadir el outfit al closet
		err = h.closetRepository.AddOutfitToCloset(context, closetID, savedOutfit.ID)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "partial_success",
				"message": "Outfit saved but not added to closet",
				"data": fiber.Map{
					"outfit":   savedOutfit,
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

func NewOutfitHandler(router fiber.Router, repository models.OutfitRepository, garmentRepository models.GarmentRepository, closetRepository models.ClosetRepository) {
	outfitGenerator := services.NewOutfitGeneratorService(garmentRepository, closetRepository)

	handler := &OutfitHandler{
		repository:        repository,
		garmentRepository: garmentRepository,
		closetRepository:  closetRepository,
		outfitGenerator:   outfitGenerator,
	}

	router.Post("/", handler.CreateOutfit)
	router.Patch("/:id", handler.UpdateOutfit)
	router.Patch("/:id/archive", handler.ArchiveOutfit)
	router.Get("/", handler.GetOutfitsByUser)
	router.Get("/generate/:closetId", handler.GenerateRandomOutfit)
	router.Get("/:id", handler.GetOutfit)
	router.Delete("/:id", handler.DeleteOutfit)
}
