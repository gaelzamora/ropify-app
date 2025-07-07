package services

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lucasb-eyer/go-colorful"
)

type OutfitGeneratorService struct {
	garmentRepository models.GarmentRepository
}

func NewOutfitGeneratorService(garmentRepository models.GarmentRepository) *OutfitGeneratorService {
	return &OutfitGeneratorService{
		garmentRepository: garmentRepository,
	}
}

// GenerateRandomOutfit crea un outfit aleatorio basado en colores y etiquetas compatibles
func (s *OutfitGeneratorService) GenerateRandomOutfit(ctx context.Context, userID uuid.UUID) (*models.Outfit, []*models.Garment, error) {
	// Obtener una prenda superior (top)
	tops, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "top", 10)
	if err != nil || len(tops) == 0 {
		return nil, nil, fmt.Errorf("no top garments found")
	}

	// Seleccionar un top aleatorio
	rand.Seed(time.Now().UnixNano())
	top := tops[rand.Intn(len(tops))]

	// Encontrar un bottom que combine con el top
	bottoms, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "bottom", 10)
	if err != nil || len(bottoms) == 0 {
		return nil, nil, fmt.Errorf("no bottom garments found")
	}

	// Filtrar bottoms por compatibilidad de color
	matchingBottoms := s.filterByColorMatch(top, bottoms)
	if len(matchingBottoms) == 0 {
		matchingBottoms = bottoms // Si no hay coincidencias, usar cualquiera
	}

	bottom := matchingBottoms[rand.Intn(len(matchingBottoms))]

	// Buscar calzado que combine
	shoes, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "sneakers", 10)
	var shoe *models.Garment
	if err == nil && len(shoes) > 0 {
		matchingShoes := s.filterByColorMatch(bottom, shoes)
		if len(matchingShoes) == 0 {
			matchingShoes = shoes
		}
		shoe = matchingShoes[rand.Intn(len(matchingShoes))]
	}

	// Buscar accesorios que complementen
	accessories, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "accessories", 5)
	var accessory *models.Garment
	if err == nil && len(accessories) > 0 {
		matchingAccessories := s.filterByColorMatch(top, accessories)
		if len(matchingAccessories) == 0 {
			matchingAccessories = accessories
		}
		accessory = matchingAccessories[rand.Intn(len(matchingAccessories))]
	}

	// Crear el outfit
	outfitGarments := []*models.Garment{top, bottom}
	garmentIDs := pq.StringArray{top.ID.String(), bottom.ID.String()}

	if shoe != nil {
		outfitGarments = append(outfitGarments, shoe)
		garmentIDs = append(garmentIDs, shoe.ID.String())
	}

	if accessory != nil {
		outfitGarments = append(outfitGarments, accessory)
		garmentIDs = append(garmentIDs, accessory.ID.String())
	}

	// Generar nombre descriptivo
	outfitName := fmt.Sprintf("%s outfit with %s", s.extractMainLabel(top), s.extractMainLabel(bottom))

	outfit := &models.Outfit{
		UserID:     userID,
		Name:       outfitName,
		GarmentIDs: garmentIDs,
		CreatedAt:  time.Now(),
	}

	return outfit, outfitGarments, nil
}

// filterByColorMatch selecciona prendas que combinen por color
func (s *OutfitGeneratorService) filterByColorMatch(baseGarment *models.Garment, candidates []*models.Garment) []*models.Garment {
	var matches []*models.Garment
	baseColor := baseGarment.Color

	for _, garment := range candidates {
		if s.areColorsCompatible(baseColor, garment.Color) {
			matches = append(matches, garment)
		}
	}

	// Si no encontramos coincidencias, devolvemos al menos algunas prendas
	if len(matches) < 2 && len(candidates) > 0 {
		return candidates
	}

	return matches
}

// areColorsCompatible determina si dos colores hacen juego
func (s *OutfitGeneratorService) areColorsCompatible(color1, color2 string) bool {
	// Colores iguales siempre combinan
	if color1 == color2 {
		return true
	}

	isHex1 := strings.HasPrefix(color1, "#")
	isHex2 := strings.HasPrefix(color2, "#")

	// Colores neutros combinan con todo
    neutralColors := []string{"#FFFFFF", "#F5F5F5", "#EEEEEE", "#000000", "#111111", "#222222", "#808080", "#A9A9A9"}
    for _, neutral := range neutralColors {
        if strings.EqualFold(color1, neutral) || strings.EqualFold(color2, neutral) {
            return true
        }
    }

    // Si no son códigos hex, no podemos calcular compatibilidad avanzada
    if !isHex1 || !isHex2 {
        return true // Por defecto permitir la combinación
    }

    // Convertir strings a objetos Color
    c1, err1 := colorful.Hex(color1)
    c2, err2 := colorful.Hex(color2)
    if err1 != nil || err2 != nil {
        return true // Si hay error, permitimos por defecto
    }

    // Convertir a HSL (Hue, Saturation, Lightness)
    h1, s1, l1 := c1.Hsl()
    h2, s2, l2 := c2.Hsl()

    // 1. Complementarios (colores opuestos en la rueda)
    hDiff := math.Abs(h1 - h2)
    if (hDiff > 150 && hDiff < 210) {
        return true // Colores opuestos (complementarios)
    }

    // 2. Análogos (colores adyacentes)
    if hDiff < 30 || hDiff > 330 {
        return true // Colores adyacentes (análogos)
    }

    // 3. Triádicos (120° en la rueda de color)
    if (hDiff > 110 && hDiff < 130) || (hDiff > 230 && hDiff < 250) {
        return true
    }

    // 4. Monocromáticos (mismo tono, diferente luminosidad)
    if hDiff < 10 && math.Abs(l1-l2) > 0.2 {
        return true
    }

    // 5. Regla básica de contraste para outfits
    // Permitir combinaciones donde un color es oscuro y otro claro
    if (l1 < 0.3 && l2 > 0.7) || (l2 < 0.3 && l1 > 0.7) {
        return true
    }

    // 6. Colores con baja saturación combinan mejor entre sí
    if s1 < 0.3 && s2 < 0.3 {
        return true
    }

    // Para otros casos, usar una probabilidad más baja (30%)
    // Esto mantiene cierta diversidad pero prioriza las reglas de color
    return rand.Float32() > 0.7
}

// extractMainLabel extrae una etiqueta principal para describir la prenda
func (s *OutfitGeneratorService) extractMainLabel(garment *models.Garment) string {
	if len(garment.Labels) > 0 {
		return garment.Labels[0]
	}
	return string(garment.Category)
}
