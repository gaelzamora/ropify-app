package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GarmentOptimized struct {
	ID       uuid.UUID `json:"id"`
	ImageURL string    `json:"image_url"`
}

type GarmentOptimizedArray []GarmentOptimized

func (a GarmentOptimizedArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *GarmentOptimizedArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("failed to unmarshal GarmentOptimizedArray value")
	}

	return json.Unmarshal(bytes, a)
}

type Outfit struct {
	ID        uuid.UUID             `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID             `json:"user_id" gorm:"type:uuid;not null"`
	Name      string                `json:"name" gorm:"not null"`
	Garments  GarmentOptimizedArray `json:"garments" gorm:"type:jsonb"`
	Occasion  string                `json:"occasion"`
	Archived  bool                  `json:"archived" gorm:"default:false"`
	ImageURL  string                `json:"image_url"`
	CreatedAt time.Time             `json:"created_at"`
}

type OutfitRepository interface {
	AddOutfit(ctx context.Context, outfit *Outfit) (*Outfit, error)
	UpdateOutfit(ctx context.Context, outfitID uuid.UUID, updateData map[string]interface{}) (*Outfit, error)
	DeleteOutfit(ctx context.Context, outfitIOD uuid.UUID) error
	ArchiveOutfit(ctx context.Context, outfitID uuid.UUID) error
	GetOutfitByID(ctx context.Context, outfitID uuid.UUID) (*Outfit, error)
	GetOutfitsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Outfit, error)
	CreateOutfitWithGarments(ctx context.Context, outfit *Outfit, garments GarmentOptimizedArray) (*Outfit, error)
}

type OutfitGarment struct {
	OutfitID  uuid.UUID `gorm:"primaryKey;type:uuid"`
	GarmentID uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
}

func (o *Outfit) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
