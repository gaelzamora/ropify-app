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

type AccessoryType string

const (
	AccessoryTypeHat     AccessoryType = "hat"
	AccessoryTypeBag     AccessoryType = "bag"
	AccessoryTypeJewelry AccessoryType = "jewelry" // Pueden ser múltiples
	AccessoryTypeGlasses AccessoryType = "glasses"
	AccessoryTypeOther   AccessoryType = "other"
)

// getAccessoryType clasifica un accesorio en su subcategoría
func (s *OutfitGeneratorService) getAccessoryType(garment *models.Garment) AccessoryType {
    // Solo procesar si es un accesorio
    if garment.Category != "accesories" {
        return AccessoryTypeOther
    }
    
    // Detectar gorras/sombreros
    if s.isHat(garment) {
        return AccessoryTypeHat
    }
    
    // Detectar bolsos/mochilas
    if s.isBag(garment) {
        return AccessoryTypeBag
    }
    
    // Detectar gafas/lentes
    if s.isGlasses(garment) {
        return AccessoryTypeGlasses
    }
    
    // Detectar joyería (collares, pulseras, anillos)
    if s.isJewelry(garment) {
        return AccessoryTypeJewelry
    }
    
    return AccessoryTypeOther
}

// isHat detecta si un accesorio es una gorra o sombrero
func (s *OutfitGeneratorService) isHat(garment *models.Garment) bool {
    hatKeywords := []string{"hat", "cap", "beanie", "headgear", "visor", "baseball cap", 
                            "cricket cap", "trucker hat", "fedora", "beret"}
    
    for _, label := range garment.Labels {
        lower := strings.ToLower(label)
        for _, keyword := range hatKeywords {
            if strings.Contains(lower, keyword) {
                return true
            }
        }
    }
    return false
}

// isBag detecta si un accesorio es una bolsa o mochila
func (s *OutfitGeneratorService) isBag(garment *models.Garment) bool {
    bagKeywords := []string{"bag", "handbag", "purse", "backpack", "tote", "clutch", 
                           "shoulder bag", "crossbody", "messenger", "pouch", "wallet"}
    
    for _, label := range garment.Labels {
        lower := strings.ToLower(label)
        for _, keyword := range bagKeywords {
            if strings.Contains(lower, keyword) {
                return true
            }
        }
    }
    return false
}

// isGlasses detecta si un accesorio es lentes/gafas
func (s *OutfitGeneratorService) isGlasses(garment *models.Garment) bool {
    glassesKeywords := []string{"glasses", "sunglasses", "eyewear", "spectacles", "shades"}
    
    for _, label := range garment.Labels {
        lower := strings.ToLower(label)
        for _, keyword := range glassesKeywords {
            if strings.Contains(lower, keyword) {
                return true
            }
        }
    }
    return false
}

// isJewelry detecta si un accesorio es joyería (puede haber múltiples)
func (s *OutfitGeneratorService) isJewelry(garment *models.Garment) bool {
    jewelryKeywords := []string{"necklace", "bracelet", "ring", "earring", "jewelry", 
                               "pendant", "chain", "anklet", "bangle", "jewellery"}
    
    for _, label := range garment.Labels {
        lower := strings.ToLower(label)
        for _, keyword := range jewelryKeywords {
            if strings.Contains(lower, keyword) {
                return true
            }
        }
    }
    return false
}

// Funciones auxiliares para categorizar prendas segun sus etiquetas
func (s *OutfitGeneratorService) isSandal(garment *models.Garment) bool {
	for _, label := range garment.Labels {
		lower := strings.ToLower(label)
		if strings.Contains(lower, "sandal") || strings.Contains(lower, "flip flop") ||
			strings.Contains(lower, "slipper") || strings.Contains(lower, "slide") ||
			strings.Contains(lower, "open toe") {
			return true
		}
	}
	return false
}

func (s *OutfitGeneratorService) isPant(garment *models.Garment) bool {
	if garment.Category == "bottom" {
		// Primero ver si es short/bermuda (para sandalias)
		for _, label := range garment.Labels {
			lower := strings.ToLower(label)
			if strings.Contains(lower, "short") || strings.Contains(lower, "bermuda") ||
				strings.Contains(lower, "trunk") || strings.Contains(lower, "boardshort") {
				return false // No es pantalón largo
			}
		}

		// Si no se detectó como short, asumir que es pantalón
		return true
	}
	return false
}

func (s *OutfitGeneratorService) isOuterLayer(garment *models.Garment) bool {
	for _, label := range garment.Labels {
		lower := strings.ToLower(label)
		if strings.Contains(lower, "hoodie") || strings.Contains(lower, "jacket") ||
			strings.Contains(lower, "coat") || strings.Contains(lower, "sweater") ||
			strings.Contains(lower, "cardigan") || strings.Contains(lower, "blazer") {
			return true
		}
	}
	return false
}

func (s *OutfitGeneratorService) validateOutfit(garments []*models.Garment) bool {
	hasOuterLayer := false
	hasInnerTop := false

	for _, g := range garments {
		if g.Category == "top" {
			hasOuterLayer = true
		} else {
			hasInnerTop = true
		}
	}

	if hasOuterLayer && !hasInnerTop {
		return false
	}

	return true
}

// GenerateRandomOutfit crea un outfit aleatorio basado en colores y etiquetas compatibles
func (s *OutfitGeneratorService) GenerateRandomOutfit(ctx context.Context, userID uuid.UUID) (*models.Outfit, []*models.Garment, error) {
	// 1. PRIMERA DECISIÓN: Elegir entre outfit basado en vestido o en conjunto top-bottom
	useDress := rand.Intn(100) < 40 // 100% de probabilidad de usar vestido

	var baseGarments []*models.Garment // Prendas base del outfit
	var garmentIDs pq.StringArray
	var outfitType string

	if useDress {
		// === CAMINO DEL VESTIDO ===
		dresses, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "dress", 10)
		if err != nil || len(dresses) == 0 {
			// Si no hay vestidos disponibles, cambiamos a conjunto top-bottom
			fmt.Println("Error al obtener vestidos: %v\n", err)
			useDress = false
		} else if len(dresses) == 0 {
			fmt.Println("No hay veestidos disponibles para el usuario")
			useDress = false
		} else {
			// Usar un vestido como base
			dress := dresses[rand.Intn(len(dresses))]
			baseGarments = append(baseGarments, dress)
			garmentIDs = append(garmentIDs, dress.ID.String())
			outfitType = "dress"
		}
	}

	if !useDress {
		// === CAMINO TOP + BOTTOM ===
		// 1. Seleccionar top
		tops, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "top", 10)
		if err != nil || len(tops) == 0 {
			return nil, nil, fmt.Errorf("no top garments found")
		}

		// Filtrar si es capa exterior (necesita capa interior)
		var innerTops, outerTops []*models.Garment
		for _, top := range tops {
			if s.isOuterLayer(top) {
				outerTops = append(outerTops, top)
			} else {
				innerTops = append(innerTops, top)
			}
		}

		// Seleccionar top básico primero
		if len(innerTops) == 0 {
			return nil, nil, fmt.Errorf("no basic tops found")
		}

		mainTop := innerTops[rand.Intn(len(innerTops))]
		baseGarments = append(baseGarments, mainTop)
		garmentIDs = append(garmentIDs, mainTop.ID.String())

		// Añadir capa exterior aleatoriamente (30% de probabilidad)
		if len(outerTops) > 0 && rand.Intn(100) < 30 {
			// Filtrar outer layers que combinen con el top
			matchingOuters := s.filterByColorMatch(mainTop, outerTops)
			if len(matchingOuters) > 0 {
				outer := matchingOuters[rand.Intn(len(matchingOuters))]
				baseGarments = append(baseGarments, outer)
				garmentIDs = append(garmentIDs, outer.ID.String())
			}
		}

		// 2. Seleccionar bottom
		bottoms, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "bottom", 10)
		if err != nil || len(bottoms) == 0 {
			return nil, nil, fmt.Errorf("no bottom garments found")
		}

		// Filtrar bottoms por compatibilidad de color
		matchingBottoms := s.filterByColorMatch(baseGarments[0], bottoms)
		if len(matchingBottoms) == 0 {
			matchingBottoms = bottoms
		}

		bottom := matchingBottoms[rand.Intn(len(matchingBottoms))]
		baseGarments = append(baseGarments, bottom)
		garmentIDs = append(garmentIDs, bottom.ID.String())
		outfitType = "casual"
	}

	// 3. CALZADO INTELIGENTE
	shoes, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "sneakers", 10)
	if err == nil && len(shoes) > 0 {
		// Filtrar sandalias si hay pantalones
		hasPants := false
		for _, garment := range baseGarments {
			if s.isPant(garment) {
				hasPants = true
				break
			}
		}

		var compatibleShoes []*models.Garment
		for _, shoe := range shoes {
			// Si hay pantalones, no usar sandalias
			if hasPants && s.isSandal(shoe) {
				continue
			}
			compatibleShoes = append(compatibleShoes, shoe)
		}

		if len(compatibleShoes) > 0 {
			// Buscar zapatos que combinen con el bottom o el vestido
			baseForShoe := baseGarments[len(baseGarments)-1] // Último elemento (bottom o dress)
			matchingShoes := s.filterByColorMatch(baseForShoe, compatibleShoes)
			if len(matchingShoes) == 0 {
				matchingShoes = compatibleShoes
			}

			shoe := matchingShoes[rand.Intn(len(matchingShoes))]
			baseGarments = append(baseGarments, shoe)
			garmentIDs = append(garmentIDs, shoe.ID.String())
		}
	}

	// 4. MÚLTIPLES ACCESORIOS
	accessories, err := s.garmentRepository.GetGarmentsByCategory(ctx, userID, "accesories", 20) // Aumentar a 20 para tener más variedad
	if err == nil && len(accessories) > 0 {
		// Mapa para rastrear los tipos de accesorios ya seleccionados
		selectedAccessoryTypes := make(map[AccessoryType]int)
		
		// Límites por tipo de accesorio
		accessoryLimits := map[AccessoryType]int{
			AccessoryTypeHat:     1, // Máximo 1 gorra/sombrero
			AccessoryTypeBag:     1, // Máximo 1 bolsa/mochila
			AccessoryTypeGlasses: 1, // Máximo 1 par de gafas
			AccessoryTypeJewelry: 3, // Máximo 3 piezas de joyería
			AccessoryTypeOther:   2, // Máximo 2 de otros tipos
		}
		
		// Número máximo total de accesorios (3-5)
		maxAccessories := rand.Intn(3) + 3
		
		// Filtrar accesorios por compatibilidad de color
		topGarment := baseGarments[0]
		matchingAccessories := s.filterByColorMatch(topGarment, accessories)
		if len(matchingAccessories) < 3 {
			matchingAccessories = accessories // Si hay pocos matches, usar todos
		}
		
		// Mezclar aleatoriamente para no favorecer siempre los mismos accesorios
		rand.Shuffle(len(matchingAccessories), func(i, j int) {
			matchingAccessories[i], matchingAccessories[j] = matchingAccessories[j], matchingAccessories[i]
		})
		
		// Añadir accesorios respetando los límites por tipo
		accessoriesAdded := 0
		for _, acc := range matchingAccessories {
			// Clasificar el accesorio
			accType := s.getAccessoryType(acc)
			
			// Verificar si ya alcanzamos el límite para este tipo
			if selectedAccessoryTypes[accType] >= accessoryLimits[accType] {
				continue
			}
			
			// Añadir el accesorio
			baseGarments = append(baseGarments, acc)
			garmentIDs = append(garmentIDs, acc.ID.String())
			selectedAccessoryTypes[accType]++
			accessoriesAdded++
			
			// Detener si llegamos al máximo total de accesorios
			if accessoriesAdded >= maxAccessories {
				break
			}
		}
	}

	// 5. NOMBRE DESCRIPTIVO MEJORADO
	var outfitName string
	if outfitType == "dress" {
		outfitName = fmt.Sprintf("%s outfit with accessories", s.extractMainLabel(baseGarments[0]))
	} else {
		// Para conjuntos top-bottom
		topLabel := s.extractMainLabel(baseGarments[0])
		bottomLabel := s.extractMainLabel(baseGarments[len(baseGarments)/2]) // Aproximadamente el bottom
		outfitName = fmt.Sprintf("%s with %s outfit", topLabel, bottomLabel)

		// Añadir estilo si es evidente
		if strings.Contains(strings.ToLower(topLabel), "formal") ||
			strings.Contains(strings.ToLower(bottomLabel), "formal") {
			outfitName = "Formal " + outfitName
		} else if len(baseGarments) > 4 {
			outfitName = "Complete " + outfitName
		}
	}

	outfit := &models.Outfit{
		UserID:     userID,
		Name:       outfitName,
		GarmentIDs: garmentIDs,
		CreatedAt:  time.Now(),
	}

	return outfit, baseGarments, nil
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
	if hDiff > 150 && hDiff < 210 {
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
