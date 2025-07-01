package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
}

type ColorInfo struct {
	Name       string
	Hex        string
	Percentage float64
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
		"shirt":      "top",
		"t-shirt":    "top",
		"t shirt":    "top",
		"polo":       "top",
		"polo shirt": "top",
		"blouse":     "top",
		"pants":      "bottom",
		"jeans":      "bottom",
		"shorts":     "bottom",
		"skirt":      "bottom",
		"dress":      "dress",
		"sneakers":   "sneakers",
		"shoes":      "sneakers",
		"hat":        "accessories",
		"jacket":     "top",
		"sweater":    "top",
		"hoodie":     "top",
		"coat":       "top",
		"backpack":   "backpack",
		"bag":        "backpack",
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

			color.Name = getColorName(c.Red, c.Green, c.Blue)
			colors = append(colors, color)
		}
	}

	return &VisionResult{
		Labels:       labelTexts,
		MainCategory: mainCategory,
		Colors:       colors,
	}, nil
}

func getColorName(r, g, b float32) string {
	// Implementación básica de mapeo de colores
	if r < 120 && g > 100 && g > r+50 && g > b+20 {
		return "green"
	} else if r > 200 && g < 100 && b < 100 {
		return "red"
	} else if r < 100 && g < 100 && b > 200 {
		return "blue"
	} else if r > 200 && g > 200 && b < 100 {
		return "yellow"
	} else if r > 200 && g > 100 && b > 200 {
		return "pink"
	} else if r < 100 && g > 100 && b > 200 {
		return "purple"
	} else if r > 200 && g > 100 && b < 100 {
		return "orange"
	} else if r > 200 && g > 200 && b > 200 {
		return "white"
	} else if r < 100 && g < 100 && b < 100 {
		return "black"
	} else if r > 100 && r < 200 && g > 100 && g < 200 && b > 100 && b < 200 {
		return "gray"
	} else if r > 150 && g > 75 && b > 0 && b < 100 {
		return "brown"
	}

	return "unknown"
}
