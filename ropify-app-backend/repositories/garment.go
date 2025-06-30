package repositories

import (
	"context"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GarmentRepository struct {
	db *gorm.DB
}

func (r *GarmentRepository) AddGarment(ctx context.Context, garment *models.Garment) (*models.Garment, error) {
	if err := r.db.WithContext(ctx).Create(garment).Error; err != nil {
		return nil, err
	}
	return garment, nil
}

func (r *GarmentRepository) FindByBarcode(ctx context.Context, barcode string) (*models.Garment, error) {
	var garment models.Garment
	if err := r.db.WithContext(ctx).Where("barcode = ?", barcode).First(&garment).Error; err != nil {
		return nil, err
	}

	return &garment, nil
}

func (r *GarmentRepository) UpdateGarment(ctx context.Context, garmentID uuid.UUID, updatedData map[string]interface{}) (*models.Garment, error) {
	var garment models.Garment
	if err := r.db.WithContext(ctx).First(&garment, "id = ?", garmentID).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&garment).Updates(updatedData).Error; err != nil {
		return nil, err
	}
	return &garment, nil
}

func (r *GarmentRepository) DeleteGarment(ctx context.Context, garmentID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Garment{}, "id = ?", garmentID).Error
}

func (r *GarmentRepository) FilterGarments(ctx context.Context, userID uuid.UUID, filters map[string]interface{}, sortBy string, limit, offset int) ([]*models.Garment, error) {
	garments := []*models.Garment{}

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	if sortBy != "" {
		query = query.Order(sortBy + " DESC")
	}

	query = query.Offset(offset).Limit(limit)

	res := query.Find(&garments)
	if res.Error != nil {
		return nil, res.Error
	}

	return garments, nil
}

func (r *GarmentRepository) UpdateGarmentImage(userId uuid.UUID, imageURL string, garmentId uuid.UUID) error {
	return r.db.Model(&models.Garment{}).
		Where("id = ? AND user_id = ?", garmentId, userId).
		Update("image_url", imageURL).Error
}

func NewGarmentRepository(db *gorm.DB) models.GarmentRepository {
	return &GarmentRepository{
		db: db,
	}
}
