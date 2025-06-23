package repositories

import (
	"context"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutfitRepository struct {
	db *gorm.DB
}

// Crear outfit
func (r *OutfitRepository) AddOutfit(ctx context.Context, outfit *models.Outfit) (*models.Outfit, error) {
	if err := r.db.WithContext(ctx).Create(outfit).Error; err != nil {
		return nil, err
	}
	return outfit, nil
}

// Editar outfit
func (r *OutfitRepository) UpdateOutfit(ctx context.Context, outfitID uuid.UUID, updateData map[string]interface{}) (*models.Outfit, error) {
	var outfit models.Outfit
	if err := r.db.WithContext(ctx).First(&outfit, "id = ?", outfitID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&outfit).Updates(updateData).Error; err != nil {
		return nil, err
	}
	return &outfit, nil
}

// Eliminar outfit
func (r *OutfitRepository) DeleteOutfit(ctx context.Context, outfitID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Outfit{}, "id = ?", outfitID).Error
}

// Archivar outfit (soft delete, ejemplo: usando un campo "archived")
func (r *OutfitRepository) ArchiveOutfit(ctx context.Context, outfitID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Outfit{}).Where("id = ?", outfitID).Update("archived", true).Error
}

// Visualizar outfit por ID
func (r *OutfitRepository) GetOutfitByID(ctx context.Context, outfitID uuid.UUID) (*models.Outfit, error) {
	var outfit models.Outfit
	if err := r.db.WithContext(ctx).First(&outfit, "id = ?", outfitID).Error; err != nil {
		return nil, err
	}
	return &outfit, nil
}

// Listar outfits de un usuario
func (r *OutfitRepository) GetOutfitsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Outfit, error) {
	outfits := []*models.Outfit{}
	res := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&outfits)
	if res.Error != nil {
		return nil, res.Error
	}
	return outfits, nil
}

func NewOutfitRepository(db *gorm.DB) models.OutfitRepository {
	return &OutfitRepository{
		db: db,
	}
}
