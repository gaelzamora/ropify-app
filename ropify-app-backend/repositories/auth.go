package repositories

import (
	"context"

	"github.com/gaelzamora/ropify-app/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func (r *AuthRepository) RegisterUser(ctx context.Context, registerData *models.AuthCredentials) (*models.User, error) {
	user := &models.User{
		FirstName: registerData.FirstName,
		LastName:  registerData.LastName,
		Email:     registerData.Email,
		Username:  registerData.Username,
		Password:  registerData.Password,
	}

	res := r.db.Model(&models.User{}).Create(user)

	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *AuthRepository) RegisterOAuthUser(ctx context.Context, user *models.User) (*models.User, error) {
	randomPass := uuid.NewString()
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPass), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user.Password = string(hashedPassword)
    
    // Usar una transacción para garantizar la consistencia
    tx := r.db.WithContext(ctx).Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }
    
    if err := tx.Create(user).Error; err != nil {
        tx.Rollback()
        return nil, err
    }
    
    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        return nil, err
    }
    
    return user, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, query interface{}, args ...interface{}) (*models.User, error) {
	user := &models.User{}

	if res := r.db.Model(user).Where(query, args...).First(user); res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func NewAuthRepository(db *gorm.DB) models.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}
