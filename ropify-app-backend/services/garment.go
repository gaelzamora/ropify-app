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
			{Type: visionpb.Feature_LABEL_DETECTION, MaxResults: 10},
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

	labels, err := client.DetectLabels(ctx, img, nil, 10)
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
		"t-shirt":      "top",
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

		// Bottoms
		"jean":          "bottom",
		"pants":         "bottom",
		"denim":         "bottom",
		"jeans":         "bottom",
		"shorts":        "bottom",
		"skirt":         "bottom",
		"trousers":      "bottom",
		"leggings":      "bottom",
		"sweatpants":    "bottom",
		"jogging pants": "bottom",
		"chinos":        "bottom",

		// Dresses
		"dress":    "dress",
		"gown":     "dress",
		"sundress": "dress",

		// Footwear
		"sneakers": "sneakers",
		"shoes":    "sneakers",
		"boots":    "sneakers",
		"sandals":  "sneakers",
		"footwear": "sneakers",

		// Accessories
		"hat":    "accessories",
		"cap":    "accessories",
		"scarf":  "accessories",
		"gloves": "accessories",
		"socks":  "accessories",
		"belt":   "accessories",

		// Bags
		"backpack": "backpack",
		"bag":      "backpack",
		"handbag":  "backpack",
		"tote":     "backpack",
		"purse":    "backpack",
		"duffel":   "backpack",
	}

	fmt.Println("========= GOOGLE VISION CATEGORIZACION ==========")
	fmt.Println("Etiquetas detectadas: ")
	for i, label := range labels {
		fmt.Printf("  %d, %s (score: %.2f)\n", i+1, label.Description, label.Score)
	}

	for _, label := range labels {
		labelTexts = append(labelTexts, label.Description)

		normLabel := strings.ToLower(strings.TrimSpace(label.Description))

		if category, ok := clothingCategories[normLabel]; ok {
			mainCategory = category
			break
		}

		for key, value := range clothingCategories {
			if strings.Contains(normLabel, key) {
				mainCategory = value
				break
			}
		}
	}

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
