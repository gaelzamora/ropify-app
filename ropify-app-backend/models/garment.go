package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	Unknown    GarmentCategory = "unknown"
)

type StringArray []string

// Scan implementa la interfaz sql.Scanner
func (sa *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal StringArray value")
	}

	return json.Unmarshal(bytes, sa)
}

// Value implementa la interfaz driver.Valuer
func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}

	return json.Marshal(sa)
}

type Garment struct {
	ID         uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID       `json:"user_id" gorm:"type:uuid;not null"`
	Category   GarmentCategory `json:"category" gorm:"not null"`
	Color      string          `json:"color" gorm:"not null"`
	Labels     StringArray     `json:"labels" gorm:"type:jsonb"`
	ImageURL   string          `json:"image_url"`
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
