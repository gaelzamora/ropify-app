package handlers

import (
	"context"
	stdCtx "context"
	"fmt"
	"io"
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

	// Obtener datos del form-data
	name := ctx.FormValue("name")
	if name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Closet name is required",
		})
	}

	// Crear un nuevo objeto closet con los datos del formulario
	closet := &models.Closet{
		ID:     uuid.New(),
		Name:   name,
		UserID: userID,
		// Puedes añadir más campos según sea necesario
	}

	// Obtener el archivo de imagen
	file, err := ctx.FormFile("closet_image")

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Failed to upload file",
		})
	}

	// Guardar el closet
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("closets/%s/%s", userID.String(), closet.ID.String())

	src, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Failed to open file",
		})
	}
	defer src.Close()
	imageBytes, err := io.ReadAll(src)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Failed to read file",
		})
	}

	imageURL, err := services.UploadToS3Bytes(imageBytes, key, file.Filename)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Failed to upload file to S3",
		})
	}
	closet.ImageURL = imageURL

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

	closetPrefix := fmt.Sprintf("closets/%s/%s/", userID.String(), closetID.String())
	if err := services.DeleteS3Folder(closetPrefix); err != nil {
		fmt.Printf("Error deleting closet folder from S3: %v\n", err)
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

// RemoveMultipleGarmentsFromCloset elimina varias prendas de un closet
func (h *ClosetHandler) RemoveMultipleGarmentsFromCloset(ctx *fiber.Ctx) error {
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

	// Estructura para el cuerpo de la petición
	var payload struct {
		GarmentIDs []string `json:"garment_ids"`
	}

	// Parsear el cuerpo de la petición
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request format",
		})
	}

	// Verificar que hay IDs para eliminar
	if len(payload.GarmentIDs) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No garment IDs provided",
		})
	}

	// Verificar que el closet exista y pertenezca al usuario
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	// Variables para seguimiento de éxitos y errores
	successCount := 0
	failedIDs := make(map[string]string)
	imageDeleteErrors := 0

	// Procesar cada ID
	for _, idStr := range payload.GarmentIDs {
		// Convertir string a UUID
		garmentID, err := uuid.Parse(idStr)
		if err != nil {
			failedIDs[idStr] = "Invalid UUID format"
			continue
		}

		imageURL, err := h.garmentRepository.GetGarmentImageURL(context, garmentID)
		if err == nil && imageURL != "" {
			// Extraer clave de S3 de la URL
			s3Key, err := services.ExtractS3KeyFromURL(imageURL)
			if err == nil {
				// Intentar eliminar la imagen de S3
				err = services.DeleteFromS3(s3Key)
				if err != nil {
					// Registrar el error pero continuar con la eliminación de la prenda
					fmt.Printf("Error deleting image from S3: %v\n", err)
					imageDeleteErrors++
				}
			}
		}

		// Intentar eliminar la prenda del closet
		err = h.closetRepository.RemoveGarmentFromCloset(context, closetID, garmentID)
		if err != nil {
			failedIDs[idStr] = err.Error()
			continue
		}

		err = h.garmentRepository.DeleteGarment(context, garmentID)
		if err != nil {
			failedIDs[idStr] = err.Error()
			continue
		}

		successCount++
	}

	// Preparar respuesta
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"removed":      successCount,
		"total":        len(payload.GarmentIDs),
		"failed_items": failedIDs,
	})
}

// FilterClosetGarments permite filtrar prendas de un closet con múltiples criterios
func (h *ClosetHandler) FilterClosetGarments(ctx *fiber.Ctx) error {
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
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	closet, err := h.closetRepository.GetClosetByID(dbCtx, closetID)
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

	// Parámetros de filtrado y paginación
	pageParam := ctx.Query("page", "1")
	limitParam := ctx.Query("limit", "50")
	category := ctx.Query("category", "")
	sortBy := ctx.Query("sort_by", "created_at")

	page, _ := strconv.Atoi(pageParam)
	limit, _ := strconv.Atoi(limitParam)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	offset := (page - 1) * limit

	// Construir mapa de filtros
	filters := make(map[string]interface{})
	if category != "" {
		filters["category"] = category
	}

	// Obtener garments con filtros
	garments, err := h.closetRepository.FilterGarmentsByCloset(dbCtx, closetID, filters, sortBy, limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Preparar respuesta
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"closet_name": closet.Name,
			"garments":    garments,
		},
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
		},
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

	closetID, err := uuid.Parse(ctx.Params("closetId"))
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

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	type GarmentImage struct {
		ID       uuid.UUID `json:"id"`
		ImageURL string    `json:"image_url"`
	}

	type OutfitWithGarments struct {
		models.Outfit
		Garments []GarmentImage `json:"garments"`
	}

	var expandedOutfits []OutfitWithGarments

	// for _, outfit := range outfits {
	// 	var garments []GarmentImage
	// 	for _, garmentID := range outfit. {
	// 		id, err := uuid.Parse(garmentID)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		garment, err := h.garmentRepository.GetGarmentByID(context, id)
	// 		if err == nil {
	// 			garments = append(garments, GarmentImage{
	// 				ID:       garment.ID,
	// 				ImageURL: garment.ImageURL,
	// 			})
	// 		}
	// 	}
	// 	expandedOutfits = append(expandedOutfits, OutfitWithGarments{
	// 		Outfit:   *outfit,
	// 		Garments: garments,
	// 	})
	// }

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   expandedOutfits,
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

// AnalyzeAndCreateGarmentInCloset analiza una imagen, crea un garment y lo añade directamente a un closet
func (h *ClosetHandler) AnalyzeAndCreateGarmentInCloset(ctx *fiber.Ctx) error {
	// Obtener y validar ID de usuario
	userIDStr := ctx.Locals("userId").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	// Obtener y validar ID del closet
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

	// Procesar el archivo de imagen
	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "No image provided",
		})
	}

	fileContent, err := file.Open()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error opening image",
		})
	}
	defer fileContent.Close()

	imageBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error reading image",
		})
	}

	// Analizar la imagen
	visionResult, err := services.AnalyzeGarmentImage(imageBytes)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error analyzing image: " + err.Error(),
		})
	}

	// Intentar remover el fondo
	imageBytesNoBg, err := services.RemoveBackground(imageBytes)
	if err != nil {
		fmt.Println("Failed to remove background, Original image will be used: ", err)
		imageBytesNoBg = imageBytes
	} else {
		fmt.Println("Image received from RemoveBackground: ", len(imageBytesNoBg))
	}

	// Subir la imagen a S3
	key := fmt.Sprintf("closets/%s/%s/garments", userID.String(), closet.ID.String())

	imageURL, err := services.UploadToS3Bytes(imageBytesNoBg, key, file.Filename)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error uploading image to S3",
		})
	}

	// Determinar la categoría en base al análisis
	var category models.GarmentCategory
	switch visionResult.MainCategory {
	case "top":
		category = models.Top
	case "bottom":
		category = models.Bottoms
	case "dress":
		category = models.Dress
	case "sneakers":
		category = models.Sneakers
	case "accessories":
		category = models.Accesories
	case "backpack":
		category = models.Backpack
	default:
		category = models.Accesories
	}

	// Determinar el color
	color := "unknown"
	if len(visionResult.Colors) > 0 {
		color = visionResult.Colors[0].Hex
	}

	// Crear el objeto garment
	garment := models.Garment{
		UserID:     userID,
		Category:   category,
		Color:      color,
		Labels:     visionResult.Labels,
		ImageURL:   imageURL,
		IsVerified: true,
		CreatedAt:  time.Now(),
	}

	dbCtx, cancelDB := stdCtx.WithTimeout(stdCtx.Background(), 10*time.Second)
	defer cancelDB()

	// Guardar el garment en la base de datos
	newGarment, err := h.garmentRepository.AddGarment(dbCtx, &garment)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error creating garment: " + err.Error(),
		})
	}

	// Añadir el garment al closet
	addCtx, addCancel := stdCtx.WithTimeout(stdCtx.Background(), 5*time.Second)
	defer addCancel()
	if err := h.closetRepository.AddGarmentToCloset(addCtx, closetID, newGarment.ID); err != nil {
		// Si falla la adición al closet, eliminar el garment para mantener consistencia
		_ = h.garmentRepository.DeleteGarment(dbCtx, newGarment.ID)

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error adding garment to closet: " + err.Error(),
		})
	}

	// Devolver respuesta exitosa
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"garment":  newGarment,
			"analysis": visionResult,
			"closet":   closet,
		},
	})
}

// NewClosetHandler crea un nuevo manejador para closets y configura la	s rutas
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
	router.Post("/:id/analyze-and-add", handler.AnalyzeAndCreateGarmentInCloset)
	router.Post("/:id/garments/remove", handler.RemoveMultipleGarmentsFromCloset)
	router.Post("/:id/garments/:garment_id", handler.AddGarmentToCloset)
	router.Get("/:id/filter-garments", handler.FilterClosetGarments)

	// Rutas para outfits en closets
	router.Post("/:id/outfits/:outfit_id", handler.AddOutfitToCloset)
	router.Delete("/:id/outfits/:outfit_id", handler.RemoveOutfitFromCloset)
	router.Get("/:closetId/outfits", handler.GetClosetOutfits)

	// Ruta para generar outfit aleatorio desde un closet
	router.Post("/:id/generate-outfit", handler.GenerateRandomOutfitFromCloset)
}
