package handlers

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/gaelzamora/ropify-app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GarmentHandler struct {
	repository models.GarmentRepository
}

func (h *GarmentHandler) AddGarment(ctx *fiber.Ctx) error {
	// 1. Parsear los datos JSON del garment
	var garment models.Garment
	if err := ctx.BodyParser(&garment); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// 2. Obtener el userId del token de autenticación
	userIdStr := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	// Asignar el userId al garment
	garment.UserID = userId

	// 3. Procesar la imagen si existe
	file, err := ctx.FormFile("garment_image")
	var imageURL string

	if err == nil && file != nil {
		// Si hay una imagen, subirla a S3
		key := fmt.Sprintf("garments/users/%s", userId.String())
		imageURL, err = services.UploadToS3(file, key)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "fail",
				"message": "Failed to upload image: " + err.Error(),
			})
		}

		// Asignar la URL de la imagen al garment
		garment.ImageURL = imageURL
	}

	// 4. Guardar el garment en la base de datos
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newGarment, err := h.repository.AddGarment(context, &garment)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   newGarment,
	})
}

func (h *GarmentHandler) FindByBarcode(ctx *fiber.Ctx) error {
	barcode := ctx.Params("barcode")
	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	garment, err := h.repository.FindByBarcode(context, barcode)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": "Garment not found",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   garment,
	})
}

func (h *GarmentHandler) DeleteGarment(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	garmentID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid garment ID",
		})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.repository.DeleteGarment(context, garmentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Garment deleted",
	})
}

func (h *GarmentHandler) UpdateGarment(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	garmentID, err := uuid.Parse(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid garment ID",
		})
	}

	var updateData map[string]interface{}
	if err := ctx.BodyParser(&updateData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error()})
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updatedGarment, err := h.repository.UpdateGarment(context, garmentID, updateData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": updatedGarment})
}

func (h *GarmentHandler) FilterGarments(ctx *fiber.Ctx) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	pageParam := ctx.Query("page", "1")
	limitParam := ctx.Query("limit", "10")
	color := ctx.Query("color", "")
	brand := ctx.Query("brand", "")
	category := ctx.Query("category", "")
	userIDParam := ctx.Query("user_id", "")

	page, _ := strconv.Atoi(pageParam)
	limit, _ := strconv.Atoi(limitParam)
	offset := (page - 1) * limit

	filters := make(map[string]interface{})
	if color != "" {
		filters["color"] = color
	}
	if brand != "" {
		filters["brand"] = brand
	}
	if category != "" {
		filters["category"] = category
	}

	var userID uuid.UUID
	if userIDParam != "" {
		userID, _ = uuid.Parse(userIDParam)
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "user_id is required",
		})
	}

	garments, err := h.repository.FilterGarments(context, userID, filters, "created_at", limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "",
		"data":    garments,
	})

}

func (h *GarmentHandler) UploadGarmentImage(ctx *fiber.Ctx) error {
	garmentIdStr := ctx.Params("id")
	garmentId, err := uuid.Parse(garmentIdStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid garment ID",
		})
	}

	file, err := ctx.FormFile("garment_image")

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Failed to upload file",
		})
	}

	userIdStr := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	key := fmt.Sprintf("garments/%s/%s", userId.String(), file.Filename)

	imageURL, err := services.UploadToS3(file, key)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Failed to upload file to S3",
		})
	}

	if err := h.repository.UpdateGarmentImage(userId, imageURL, garmentId); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "Failed to update garment image",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "Garment image uploaded successfully",
		"data":    imageURL,
	})
}

func (h *GarmentHandler) LookupByBarcode(ctx *fiber.Ctx) error {
	var payload struct{ Barcode string }

	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid request data",
		})
	}

	productData, err := services.BarcodeLookup(payload.Barcode)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error calling external API: " + err.Error(),
		})
	}

	userIdStr := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

	// Mapear la categoría de la API a tu enum de GarmentCategory
	category := models.Unknown // Valor predeterminado
	switch productData.Category {
	case "Tops", "Shirt", "T-Shirt":
		category = models.Top
	case "Bottoms", "Pants", "Jeans", "Shorts":
		category = models.Bottoms
	case "Dress", "Dresses":
		category = models.Dress
	case "Sneakers", "Shoes":
		category = models.Sneakers
	case "Accessories", "Jewelry", "Watches":
		category = models.Accesories
	case "Backpack", "Bag":
		category = models.Backpack
	default:
		category = models.Unknown // Por defecto
	}

	garment := models.Garment{
		ID:         uuid.New(),
		UserID:     userId,
		Name:       productData.ProductName,
		Category:   category,
		Color:      productData.Color,
		Brand:      productData.Brand,
		Size:       productData.Size,
		ImageURL:   productData.ImageURL,
		Barcode:    productData.Barcode,
		IsVerified: true, // Asumimos que los productos de la API son verificados
		CreatedAt:  time.Now(),
	}

	ctxDB, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newGarment, err := h.repository.AddGarment(ctxDB, &garment)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error saving garment in DB: " + err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status": "success",
		"data":   newGarment,
	})
}

func (h *GarmentHandler) AnalyzeAndCreateGarment(ctx *fiber.Ctx) error {
	userIdStr := ctx.Locals("userId").(string)
	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user ID",
		})
	}

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

	visionResult, err := services.AnalyzeGarmentImage(imageBytes)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error analying image: " + err.Error(),
		})
	}

	key := fmt.Sprintf("garments/users/%s", userId.String())
	imageURL, err := services.UploadToS3(file, key)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error uploading image to S3",
		})
	}

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
		category = models.Unknown
	}

	color := "unknown"
	if len(visionResult.Colors) > 0 {
		color = visionResult.Colors[0].Name
	}

	garment := models.Garment{
		UserID:     userId,
		Name:       fmt.Sprintf("%s %s", color, visionResult.MainCategory),
		Category:   category,
		Color:      color,
		Brand:      "", // Podría ser completado por el usuario después
		Size:       "", // Podría ser completado por el usuario después
		ImageURL:   imageURL,
		IsVerified: true,
		CreatedAt:  time.Now(),
	}

	context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newGarment, err := h.repository.AddGarment(context, &garment)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Error al crear prenda: " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"garment":  newGarment,
			"analysis": visionResult,
		},
	})
}

func NewGarmentHandler(router fiber.Router, repository models.GarmentRepository) {
	handler := &GarmentHandler{
		repository: repository,
	}

	router.Post("/", handler.AddGarment)
	router.Post("/barcode", handler.LookupByBarcode)
	router.Post("/analyze", handler.AnalyzeAndCreateGarment)
	router.Post("/:id", handler.UploadGarmentImage)

	router.Get("/", handler.FilterGarments)
	router.Get("/:barcode", handler.FindByBarcode)

	router.Patch("/:id", handler.UpdateGarment)

	router.Delete("/:id", handler.DeleteGarment)
}
