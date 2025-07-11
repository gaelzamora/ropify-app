package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Closet struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Name      string         `json:"name" gorm:"not null"`
	Garments  []Garment      `json:"garments,omitempty" gorm:"many2many:closet_garments;"`
	OutfitIDs pq.StringArray `json:"outfit_ids" gorm:"type:uuid[];column:outfit_ids"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ClosetGarment struct {
	ClosetID  uuid.UUID `gorm:"primaryKey;type:uuid"`
	GarmentID uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
}

// ClosetRepository define métodos para trabajar con closets
type ClosetRepository interface {
	AddCloset(ctx context.Context, closet *Closet) (*Closet, error)
	GetClosetByID(ctx context.Context, id uuid.UUID) (*Closet, error)
	GetClosetsByUserID(ctx context.Context, userID uuid.UUID) ([]*Closet, error)
	UpdateCloset(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteCloset(ctx context.Context, id uuid.UUID) error

	// Nuevos métodos para manejar la relación con garments
	AddGarmentToCloset(ctx context.Context, closetID, garmentID uuid.UUID) error
	RemoveGarmentFromCloset(ctx context.Context, closetID, garmentID uuid.UUID) error
	GetGarmentsByCloset(ctx context.Context, closetID uuid.UUID) ([]*Garment, error)

	// Nuevos métodos para manejar la relación con outfits
	AddOutfitToCloset(ctx context.Context, closetID, outfitID uuid.UUID) error
	RemoveOutfitFromCloset(ctx context.Context, closetID, outfitID uuid.UUID) error
	GetOutfitsByCloset(ctx context.Context, closetID uuid.UUID) ([]*Outfit, error)

	GetGarmentsByCategoryAndCloset(ctx context.Context, closetID uuid.UUID,
		category string, limit int) ([]*Garment, error)
}
