package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type OutfitHandler struct {
	repository models.OutfitRepository
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

func NewOutfitHandler(router fiber.Router, repository models.OutfitRepository) {
	handler := &OutfitHandler{
		repository: repository,
	}

	router.Post("/", handler.CreateOutfit)
	router.Patch("/:id", handler.UpdateOutfit)
	router.Patch("/:id/archive", handler.ArchiveOutfit)
	router.Get("/", handler.GetOutfitsByUser)
	router.Get("/:id", handler.GetOutfit)
	router.Delete("/:id", handler.DeleteOutfit)	
}
