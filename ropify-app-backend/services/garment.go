package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

// BarcodeAPIResponse represents the response structure for a barcode lookup.
type BarcodeAPIResponse struct {
	ProductName string  `json:"product_name"`
	Brand       string  `json:"brand"`
	Category    string  `json:"category"`
	Color       string  `json:"color"`
	Barcode     *string `json:"barcode"`
	Size        string  `json:"size"`
	ImageURL    string  `json:"imageURL"`
}

type VisionResult struct {
	Labels       []string
	MainCategory string
	Colors       []ColorInfo

	ObjectMask   []byte
	BoundingPoly []Point
}

type ColorInfo struct {
	Hex        string
	Percentage float64
}

type Point struct {
	X float64
	Y float64
}

func RemoveBackground(imageBytes []byte) ([]byte, error) {
	url := "http://background-removal-service:8000/remove-background"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", "garment.png")
	if err != nil {
		return nil, fmt.Errorf("Failed to create form file: %v", err)
	}
	part.Write(imageBytes)
	writer.Close()

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to remove background: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to remove background code: %d", resp.StatusCode)
	}

	output, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read output image: %v", err)
	}

	return output, nil
}

func BarcodeLookup(barcode string) (*BarcodeAPIResponse, error) {
	client := &http.Client{Timeout: time.Second * 10}

	apiKey := os.Getenv("BARCODELOOKUP_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("BARCODELOOKUP_API_KEY environment variable not set")
	}

	// Endpoint for Barcode Lookup API
	url := fmt.Sprintf("https://api.barcodelookup.com/v3/products?barcode=%s&key=%s", barcode, apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// You can add more detailed error handling here by parsing the error response body
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	// This struct matches the response from Barcode Lookup API
	var result struct {
		Products []struct {
			ProductName string   `json:"title"`
			Brand       string   `json:"brand"`
			Category    string   `json:"category"`
			Images      []string `json:"images"`
			Color       string   `json:"color"`
			Size        string   `json:"size"`
		} `json:"products"`
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	fmt.Println("Raw API Response: ", string(bodyBytes))

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Products) > 0 {
		item := result.Products[0]

		barcodeStr := barcode
		imageURL := ""
		if len(item.Images) > 0 {
			imageURL = item.Images[0]
		}

		// Map the data from the new API to our consistent BarcodeAPIResponse struct
		return &BarcodeAPIResponse{
			ProductName: item.ProductName,
			Brand:       item.Brand,
			Category:    item.Category,
			Color:       item.Color,
			Size:        item.Size,
			Barcode:     &barcodeStr,
			ImageURL:    imageURL,
		}, nil
	}

	return nil, fmt.Errorf("no items found for that barcode")
}

func AnalyzeGarmentImage(imageData []byte) (*VisionResult, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)

	if err != nil {
		return nil, fmt.Errorf("Error al crear cliente de Vision API: %v", err)
	}

	defer client.Close()

	img := &visionpb.Image{
		Content: imageData,
	}

	// Añadir solicitud de localización de objetos
	req := &visionpb.AnnotateImageRequest{
		Image: img,
		Features: []*visionpb.Feature{
			{Type: visionpb.Feature_LABEL_DETECTION, MaxResults: 20},
			{Type: visionpb.Feature_IMAGE_PROPERTIES, MaxResults: 10},
			{Type: visionpb.Feature_OBJECT_LOCALIZATION, MaxResults: 5}, // Añadido
		},
	}

	resp, err := client.BatchAnnotateImages(ctx, &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	})
	if err != nil {
		return nil, fmt.Errorf("Error en la llamada a Vision API: %v", err)
	}

	labels, err := client.DetectLabels(ctx, img, nil, 20)
	if err != nil {
		return nil, fmt.Errorf("error al detectar etiquetas: %v", err)
	}

	props, err := client.DetectImageProperties(ctx, img, nil)
	if err != nil {
		return nil, fmt.Errorf("Error al detectar propiedades: %v", err)
	}

	var labelTexts []string
	var mainCategory string

	clothingCategories := map[string]string{
		"shirt":        "top",
		"active shirt": "top",
		"t shirt":      "top",
		"polo":         "top",
		"polo shirt":   "top",
		"blouse":       "top",
		"jacket":       "top",
		"sweater":      "top",
		"hoodie":       "top",
		"coat":         "top",
		"sweatshirt":   "top",
		"jersey":       "top",
		"cardigan":     "top",
		"button shirt": "top",
		"long sleeve":  "top", // Detecta prendas de manga larga
		"sleeve":       "top",

		// Bottoms
		"jean":           "bottom",
		"pants":          "bottom",
		"denim":          "bottom",
		"jeans":          "bottom",
		"shorts":         "bottom",
		"bermuda shorts": "bottom",
		"active shorts":  "bottom",
		"skirt":          "bottom",
		"trousers":       "bottom",
		"leggings":       "bottom",
		"sweatpants":     "bottom",
		"jogging pants":  "bottom",
		"chinos":         "bottom",

		// Dresses
		"dress":             "dress",
		"gown":              "dress",
		"sundress":          "dress",
		"cocktail dress":    "dress",
		"day dress":         "dress",
		"formal wear":       "dress",
		"one-piece garment": "dress",
		"evening gown":      "dress",
		"haute couture":     "dress",

		// Sneakers
		"sneakers":         "sneakers",
		"shoes":            "sneakers",
		"boots":            "sneakers",
		"sandals":          "sneakers",
		"footwear":         "sneakers",
		"athletic shoes":   "sneakers",
		"sports shoes":     "sneakers",
		"basketball shoes": "sneakers",
		"running shoes":    "sneakers",
		"trainers":         "sneakers",
		"jordan":           "sneakers", // Específicamente para Jordan
		"nike":             "sneakers", // Marcas comunes
		"adidas":           "sneakers",
		"puma":             "sneakers",
		"converse":         "sneakers",
		"air jordan":       "sneakers",
		"high tops":        "sneakers",
		"sneaker":          "sneakers", // Singular
		"walking shoe":     "sneakers",
		"skate shoe":       "sneakers",
		"outdoor shoe":     "sneakers",
		"basketball shoe":  "sneakers", // Singular (ya tienes el plural)
		"shoe":             "sneakers", // Genérico
		"athletic shoe":    "sneakers", // Singular
		"slipper":          "sneakers",

		// Accessories
		"hat":          "accessories",
		"cap":          "accessories",
		"scarf":        "accessories",
		"gloves":       "accessories",
		"socks":        "accessories",
		"belt":         "accessories",
		"glasses":      "accesories",
		"jewellery":    "accesories",
		"body jewelry": "accesories",
		"bag":          "accesories",
		"handbag":      "accesories",
		"shoulder bag": "accesories",

		// Bags
		"backpack":       "backpack",
		"tote":           "backpack",
		"purse":          "backpack",
		"duffel":         "backpack",
		"luggage & bags": "backpack",
		"baggage":        "backpack",
		"pocket":         "backpack",
	}

	priorityLabels := map[string]int{
		// Prioridad alta para calzado
		"sneakers":          100,
		"jeans":             100,
		"jean":              100,
		"shoes":             95,
		"athletic shoes":    90,
		"basketball shoes":  90,
		"shoulder bag":      90,
		"running shoes":     90,
		"jordan":            100,
		"air jordan":        100,
		"slipper":           100,
		"walking shoe":      95,
		"skate shoe":        95,
		"outdoor shoe":      95,
		"basketball shoe":   95, // Singular
		"shoe":              90,
		"athletic shoe":     95,
		"nike":              85,
		"adidas":            85,
		"puma":              85,
		"converse":          85,
		"sleeve":            90,
		"active shirt":      90,
		"footwear":          80,
		"trainers":          85,
		"dress":             80,
		"cocktail dress":    80,
		"day dress":         80,
		"formal wear":       75,
		"one-piece garment": 75,
		"evening gown":      75,
		"haute couture":     70,
		"hat":               95,
		"cap":               95,
		"headgear":          98,
		"baseball cap":      98,
		"trucker hat":       95,
		"visor":             90,
		"cricket cap":       95,
		"jewellery":         90,
		"body jewelry":      90,

		// Otros accesorios
		"pocket":  90,
		"glasses": 85,
		"belt":    85,
		"gloves":  85,
		"scarf":   85,
		"socks":   85,
		"baggage": 85,

		// Prioridad media para tops
		"shirt":          70,
		"t-shirt":        70,
		"hoodie":         70,
		"jacket":         70,
		"sweater":        70,
		"luggage & bags": 70,

		// Prioridad media para bottoms
		"pants":          100,
		"shorts":         100,
		"skirt":          65,
		"bermuda shorts": 100,
		"active shorts":  90,

		// Prioridad baja para materiales (ya que pueden aparecer en cualquier tipo de prenda)
		"denim":   40,
		"cotton":  30,
		"leather": 30,
		"fabric":  20,
		"textile": 10,
	}

	highestPriority := -1
	bestCategory := ""

	singularize := func(word string) string {
		// Reglas básicas para convertir a singular
		if strings.HasSuffix(word, "shoes") {
			return strings.TrimSuffix(word, "s")
		}
		if strings.HasSuffix(word, "s") && len(word) > 3 {
			// Solo quitar la 's' si la palabra es larga para evitar casos como "is" -> "i"
			return strings.TrimSuffix(word, "s")
		}
		return word
	}

	findCategoryForCompoundTerm := func(term string) (string, int) {
		// Buscar coincidencia exacta primero
		if category, ok := clothingCategories[term]; ok {
			priority := 50
			if p, exists := priorityLabels[term]; exists {
				priority = p
			}
			return category, priority
		}

		// Luego buscar términos compuestos que contengan "dress"
		if strings.Contains(term, "dress") {
			return "dress", 75 // Alta prioridad para cualquier término que contenga "dress"
		}

		return "", -1 // No encontrado
	}

	// En tu bucle de labels:
	for _, label := range labels {
		labelTexts = append(labelTexts, label.Description)
		normLabel := strings.ToLower(strings.TrimSpace(label.Description))
		singularLabel := singularize(normLabel) // Normaliza a singular también

		// Buscar categoría para el término compuesto
		if category, priority := findCategoryForCompoundTerm(normLabel); category != "" {
			if priority > highestPriority {
				bestCategory = category
				highestPriority = priority
			}
			continue // Pasar a la siguiente etiqueta
		}

		// Buscar coincidencia con ambas versiones
		for _, labelVariation := range []string{normLabel, singularLabel} {
			if category, ok := clothingCategories[labelVariation]; ok {
				priority := 50
				if p, exists := priorityLabels[labelVariation]; exists {
					priority = p
				}

				if priority > highestPriority {
					bestCategory = category
					highestPriority = priority
				}
			}
		}

		// Buscar coincidencia parcial en categorías
		for key, value := range clothingCategories {
			if strings.Contains(normLabel, key) {
				// Para accesorios, damos mayor peso a coincidencias exactas
				if value == "accessories" && normLabel == key {
					priority := 90 // Alta prioridad para coincidencias exactas de accesorios
					if p, exists := priorityLabels[key]; exists {
						priority = p
					}

					if priority > highestPriority {
						bestCategory = value
						highestPriority = priority
					}
					continue
				}

				// Comprobar prioridad para coincidencias parciales
				priority := 40 // Prioridad por defecto más baja
				if p, exists := priorityLabels[key]; exists {
					priority = p - 5 // Ligera penalización por ser coincidencia parcial
				}

				if priority > highestPriority {
					bestCategory = value
					highestPriority = priority
				}
			}
		}
	}

	// Usar la categoría con mayor prioridad si se encontró alguna
	if bestCategory != "" {
		mainCategory = bestCategory
	}

	// Resto del código existente...

	// Para debug - imprimir todas las etiquetas detectadas con sus puntajes
	fmt.Println("Todas las etiquetas detectadas:")
	for i, label := range labels {
		fmt.Printf("%d. %s (score: %.2f)\n", i+1, label.Description, label.Score)
	}

	fmt.Printf("Categoría seleccionada: %s (prioridad: %d)\n", mainCategory, highestPriority)

	var colors []ColorInfo
	if props != nil && props.DominantColors != nil {
		for _, colorInfo := range props.DominantColors.Colors {
			c := colorInfo.Color
			color := ColorInfo{
				Hex:        fmt.Sprintf("#%02X%02X%02X", int(c.Red), int(c.Green), int(c.Blue)),
				Percentage: float64(colorInfo.Score),
			}

			colors = append(colors, color)
		}
	}

	// Procesar resultados...
	var boundingPoly []Point
	if len(resp.Responses) > 0 && len(resp.Responses[0].LocalizedObjectAnnotations) > 0 {
		// Buscar el objeto de ropa con mayor score
		bestObject := resp.Responses[0].LocalizedObjectAnnotations[0]
		for _, vertex := range bestObject.BoundingPoly.NormalizedVertices {
			boundingPoly = append(boundingPoly, Point{
				X: float64(vertex.X),
				Y: float64(vertex.Y),
			})
		}
	}

	return &VisionResult{
		Labels:       labelTexts,
		MainCategory: mainCategory,
		Colors:       colors,
		BoundingPoly: boundingPoly,
	}, nil
}
