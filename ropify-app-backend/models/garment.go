package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GarmentCategory string

const (
	Top        GarmentCategory = "top"
	Bottoms    GarmentCategory = "bottom"
	Dress      GarmentCategory = "dress"
	Sneakers   GarmentCategory = "snearkers"
	Accesories GarmentCategory = "accesories"
	Backpack   GarmentCategory = "backpack"
)

type Garment struct {
	ID         uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Name       string          `json:"name" gorm:"not null"`
	Category   GarmentCategory `json:"category" gorm:"not null"`
	Color      string          `json:"color" gorm:"not null"`
	Brand      string          `json:"brand" gorm:"not null"`
	Size       string          `json:"size" gorm:"not null"`
	ImageURL   string          `json:"image_url"`
	Barcode    *string         `json:"barcode"`
	IsVerified bool            `json:"is_verified"`
	CreatedAt  time.Time       `json:"created_at"`
}

type GarmentRepository interface {
	AddGarment(ctx context.Context, garment *Garment) (*Garment, error)
	FindByBarcode(ctx context.Context, barcode string) (*Garment, error)

	UpdateGarment(ctx context.Context, garmentID uuid.UUID, updatedData map[string]interface{}) (*Garment, error)

	DeleteGarment(ctx context.Context, garmentID uuid.UUID) error

	FilterGarments(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, sortBy string, limit, offset int) ([]*Garment, error)

	UpdateGarmentImage(userId uuid.UUID, imageURL string, garmentId uuid.UUID) error
}

func (g *Garment) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID = uuid.New()
	return
}
