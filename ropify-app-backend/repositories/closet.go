package repositories

import (
	"context"
	"time"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClosetRepository struct {
	db *gorm.DB
}

func NewClosetRepository(db *gorm.DB) *ClosetRepository {
	return &ClosetRepository{db: db}
}

// AddCloset crea un nuevo closet en la base de datos
func (r *ClosetRepository) AddCloset(ctx context.Context, closet *models.Closet) (*models.Closet, error) {
	// Asegurar que se asigna un ID si no tiene uno
	if closet.ID == uuid.Nil {
		closet.ID = uuid.New()
	}

	// Establecer las marcas de tiempo
	now := time.Now()
	closet.CreatedAt = now
	closet.UpdatedAt = now

	err := r.db.WithContext(ctx).Create(&closet).Error
	if err != nil {
		return nil, err
	}

	return closet, nil
}

// GetClosetByID obtiene un closet por su ID
func (r *ClosetRepository) GetClosetByID(ctx context.Context, id uuid.UUID) (*models.Closet, error) {
	var closet models.Closet
	if err := r.db.WithContext(ctx).First(&closet, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &closet, nil
}

// GetClosetsByUserID obtiene todos los closets de un usuario
func (r *ClosetRepository) GetClosetsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Closet, error) {
	var closets []*models.Closet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&closets).Error; err != nil {
		return nil, err
	}
	return closets, nil
}

// UpdateCloset actualiza un closet existente
func (r *ClosetRepository) UpdateCloset(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	// Asegurar que UpdatedAt se actualiza
	updates["updated_at"] = time.Now()

	return r.db.WithContext(ctx).Model(&models.Closet{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteCloset elimina un closet
func (r *ClosetRepository) DeleteCloset(ctx context.Context, id uuid.UUID) error {
	// Usar una transacción para eliminar el closet y sus relaciones
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Eliminar las relaciones con garments
	if err := tx.Where("closet_id = ?", id).Delete(&models.ClosetGarment{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Eliminar el closet
	if err := tx.Delete(&models.Closet{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// AddGarmentToCloset añade una prenda a un closet
func (r *ClosetRepository) AddGarmentToCloset(ctx context.Context, closetID, garmentID uuid.UUID) error {
	closetGarment := models.ClosetGarment{
		ClosetID:  closetID,
		GarmentID: garmentID,
		CreatedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(&closetGarment).Error
}

// RemoveGarmentFromCloset elimina una prenda de un closet
func (r *ClosetRepository) RemoveGarmentFromCloset(ctx context.Context, closetID, garmentID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("closet_id = ? AND garment_id = ?", closetID, garmentID).
		Delete(&models.ClosetGarment{}).Error
}

// GetGarmentsByCloset obtiene todas las prendas de un closet
func (r *ClosetRepository) GetGarmentsByCloset(ctx context.Context, closetID uuid.UUID) ([]*models.Garment, error) {
	var garments []*models.Garment

	err := r.db.WithContext(ctx).
		Table("garments").
		Joins("JOIN closet_garments ON garments.id = closet_garments.garment_id").
		Where("closet_garments.closet_id = ?", closetID).
		Find(&garments).Error

	return garments, err
}

// GetGarmentsByCategoryAndCloset obtiene prendas de un closet filtrando por categoría (método de utilidad adicional)
func (r *ClosetRepository) GetGarmentsByCategoryAndCloset(ctx context.Context, closetID uuid.UUID,
	category string, limit int) ([]*models.Garment, error) {

	var garments []*models.Garment
	query := r.db.WithContext(ctx).
		Table("garments").
		Joins("JOIN closet_garments ON garments.id = closet_garments.garment_id").
		Where("closet_garments.closet_id = ?", closetID)

	if category != "" {
		query = query.Where("garments.category = ?", category)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&garments).Error
	return garments, err
}

// AddOutfitToCloset añade un outfit a un closet
func (r *ClosetRepository) AddOutfitToCloset(ctx context.Context, closetID, outfitID uuid.UUID) error {
	// Actualiza el campo outfit_ids que es un array de UUIDs
	return r.db.WithContext(ctx).Exec(`
        UPDATE closets 
        SET outfit_ids = array_append(outfit_ids, ?), 
            updated_at = ? 
        WHERE id = ?`,
		outfitID.String(), time.Now(), closetID).Error
}

// RemoveOutfitFromCloset elimina un outfit de un closet
func (r *ClosetRepository) RemoveOutfitFromCloset(ctx context.Context, closetID, outfitID uuid.UUID) error {
	// Elimina un UUID específico del array outfit_ids
	return r.db.WithContext(ctx).Exec(`
        UPDATE closets 
        SET outfit_ids = array_remove(outfit_ids, ?),
            updated_at = ?
        WHERE id = ?`,
		outfitID.String(), time.Now(), closetID).Error
}

// GetOutfitsByCloset obtiene todos los outfits asociados a un closet
func (r *ClosetRepository) GetOutfitsByCloset(ctx context.Context, closetID uuid.UUID) ([]*models.Outfit, error) {
	// Primero obtenemos el closet para acceder al array de outfit_ids
	var closet models.Closet
	if err := r.db.WithContext(ctx).Select("outfit_ids").First(&closet, "id = ?", closetID).Error; err != nil {
		return nil, err
	}

	if len(closet.OutfitIDs) == 0 {
		return []*models.Outfit{}, nil
	}

	// Luego obtenemos todos los outfits correspondientes
	var outfits []*models.Outfit
	if err := r.db.WithContext(ctx).
		Where("id IN (?)", closet.OutfitIDs).
		Find(&outfits).Error; err != nil {
		return nil, err
	}

	return outfits, nil
}
