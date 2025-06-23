package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Outfit struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name       string         `json:"name" gorm:"not null"`
	GarmentIDs pq.StringArray `json:"garment_ids" gorm:"type:uuid[]"`
	Tags       pq.StringArray `json:"tags" gorm:"type:text[]"`
	Occasion   string         `json:"occasion"`
	Season     string         `json:"season"`
	Archived   bool           `json:"archived" gorm:"default:false"`
	ImageURL   string         `json:"image_url"`
	CreatedAt  time.Time      `json:"created_at"`
}

type OutfitRepository interface {
	AddOutfit(ctx context.Context, outfit *Outfit) (*Outfit, error)
	UpdateOutfit(ctx context.Context, outfitID uuid.UUID, updateData map[string]interface{}) (*Outfit, error)
	DeleteOutfit(ctx context.Context, outfitIOD uuid.UUID) error
	ArchiveOutfit(ctx context.Context, outfitID uuid.UUID) error
	GetOutfitByID(ctx context.Context, outfitID uuid.UUID) (*Outfit, error)
	GetOutfitsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Outfit, error)
}

func (o *Outfit) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
