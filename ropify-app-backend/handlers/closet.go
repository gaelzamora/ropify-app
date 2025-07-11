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
)

// ClosetHandler maneja las peticiones HTTP relacionadas con closets
type ClosetHandler struct {
	closetRepository  models.ClosetRepository
	garmentRepository models.GarmentRepository
	outfitRepository  models.OutfitRepository
	outfitGenerator   *services.OutfitGeneratorService
}

// CreateCloset crea un nuevo closet
func (h *ClosetHandler) CreateCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	// Parsear datos del closet
	closet := new(models.Closet)
	if err := ctx.BodyParser(closet); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request data",
		})
	}

	// Asignar el ID del usuario
	closet.UserID = userID

	// Guardar el closet
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	savedCloset, err := h.closetRepository.AddCloset(context, closet)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   savedCloset,
	})
}

// GetCloset obtiene un closet por su ID
func (h *ClosetHandler) GetCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
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

	// Verificar que el closet pertenezca al usuario
	if closet.UserID != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You don't have access to this closet",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   closet,
	})
}

// GetUserClosets obtiene todos los closets de un usuario
func (h *ClosetHandler) GetUserClosets(ctx *fiber.Ctx) error {
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

	closets, err := h.closetRepository.GetClosetsByUserID(context, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   closets,
	})
}

// UpdateCloset actualiza un closet existente
func (h *ClosetHandler) UpdateCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	// Parsear actualizaciones
	updates := make(map[string]interface{})
	if err := ctx.BodyParser(&updates); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request data",
		})
	}

	// Solo permitir actualizar ciertos campos
	allowedUpdates := map[string]bool{
		"name": true,
	}

	filteredUpdates := make(map[string]interface{})
	for k, v := range updates {
		if _, ok := allowedUpdates[k]; ok {
			filteredUpdates[k] = v
		}
	}

	if err := h.closetRepository.UpdateCloset(context, closetID, filteredUpdates); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Obtener el closet actualizado
	updatedCloset, err := h.closetRepository.GetClosetByID(context, closetID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updatedCloset,
	})
}

// DeleteCloset elimina un closet
func (h *ClosetHandler) DeleteCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	if err := h.closetRepository.DeleteCloset(context, closetID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Closet deleted successfully",
	})
}

// AddGarmentToCloset añade una prenda a un closet
func (h *ClosetHandler) AddGarmentToCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	garmentID, err := uuid.Parse(ctx.Params("garment_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid garment ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	// Verificar que la prenda exista y pertenezca al usuario
	garment, err := h.garmentRepository.GetGarmentByID(context, garmentID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Garment not found",
		})
	}

	if garment.UserID != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You don't have access to this garment",
		})
	}

	// Añadir la prenda al closet
	if err := h.closetRepository.AddGarmentToCloset(context, closetID, garmentID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Garment added to closet successfully",
	})
}

// RemoveGarmentFromCloset elimina una prenda de un closet
func (h *ClosetHandler) RemoveGarmentFromCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	garmentID, err := uuid.Parse(ctx.Params("garment_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid garment ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	// Eliminar la prenda del closet
	if err := h.closetRepository.RemoveGarmentFromCloset(context, closetID, garmentID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Garment removed from closet successfully",
	})
}

// GetClosetGarments obtiene todas las prendas de un closet
func (h *ClosetHandler) GetClosetGarments(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	// Opcionalmente filtrar por categoría
	category := ctx.Query("category", "")

	// Opcionalmente limitar resultados
	limit, _ := strconv.Atoi(ctx.Query("limit", "50"))

	var garments []*models.Garment
	var getErr error

	if category != "" {
		garments, getErr = h.closetRepository.GetGarmentsByCategoryAndCloset(context, closetID, category, limit)
	} else {
		garments, getErr = h.closetRepository.GetGarmentsByCloset(context, closetID)
	}

	if getErr != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": getErr.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   garments,
	})
}

// AddOutfitToCloset añade un outfit a un closet
func (h *ClosetHandler) AddOutfitToCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	outfitID, err := uuid.Parse(ctx.Params("outfit_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid outfit ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	// Verificar que el outfit exista y pertenezca al usuario
	outfit, err := h.outfitRepository.GetOutfitByID(context, outfitID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Outfit not found",
		})
	}

	if outfit.UserID != userID {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "fail",
			"message": "You don't have access to this outfit",
		})
	}

	// Añadir el outfit al closet
	if err := h.closetRepository.AddOutfitToCloset(context, closetID, outfitID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Outfit added to closet successfully",
	})
}

// RemoveOutfitFromCloset elimina un outfit de un closet
func (h *ClosetHandler) RemoveOutfitFromCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	outfitID, err := uuid.Parse(ctx.Params("outfit_id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid outfit ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	// Eliminar el outfit del closet
	if err := h.closetRepository.RemoveOutfitFromCloset(context, closetID, outfitID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Outfit removed from closet successfully",
	})
}

// GetClosetOutfits obtiene todos los outfits de un closet
func (h *ClosetHandler) GetClosetOutfits(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	outfits, err := h.closetRepository.GetOutfitsByCloset(context, closetID)
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

// GenerateRandomOutfitFromCloset genera un outfit aleatorio usando las prendas de un closet específico
func (h *ClosetHandler) GenerateRandomOutfitFromCloset(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	closetID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid closet ID",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
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

	// Generar outfit aleatorio desde este closet
	outfit, garments, err := h.outfitGenerator.GenerateRandomOutfitFromCloset(context, closetID)
	if err != nil {
		if strings.Contains(err.Error(), "no top garments found") ||
			strings.Contains(err.Error(), "no bottom garments found") ||
			strings.Contains(err.Error(), "no basic tops found") {

			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "success",
				"data": fiber.Map{
					"outfit":   nil,
					"garments": []*models.Garment{},
					"message":  "Not enough essential garments in this closet to generate an outfit.",
				},
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Si se solicita guardar el outfit
	if ctx.Query("save", "false") == "true" {
		savedOutfit, err := h.outfitRepository.AddOutfit(context, outfit)
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

// NewClosetHandler crea un nuevo manejador para closets y configura las rutas
func NewClosetHandler(
	router fiber.Router,
	closetRepo models.ClosetRepository,
	garmentRepo models.GarmentRepository,
	outfitRepo models.OutfitRepository,
	outfitGenerator *services.OutfitGeneratorService,
) {
	handler := &ClosetHandler{
		closetRepository:  closetRepo,
		garmentRepository: garmentRepo,
		outfitRepository:  outfitRepo,
		outfitGenerator:   outfitGenerator,
	}

	// Rutas para closets
	router.Post("/", handler.CreateCloset)
	router.Get("/", handler.GetUserClosets)
	router.Get("/:id", handler.GetCloset)
	router.Put("/:id", handler.UpdateCloset)
	router.Delete("/:id", handler.DeleteCloset)

	// Rutas para garments en closets
	router.Post("/:id/garments/:garment_id", handler.AddGarmentToCloset)
	router.Delete("/:id/garments/:garment_id", handler.RemoveGarmentFromCloset)
	router.Get("/:id/garments", handler.GetClosetGarments)

	// Rutas para outfits en closets
	router.Post("/:id/outfits/:outfit_id", handler.AddOutfitToCloset)
	router.Delete("/:id/outfits/:outfit_id", handler.RemoveOutfitFromCloset)
	router.Get("/:id/outfits", handler.GetClosetOutfits)

	// Ruta para generar outfit aleatorio desde un closet
	router.Post("/:id/generate-outfit", handler.GenerateRandomOutfitFromCloset)
}
