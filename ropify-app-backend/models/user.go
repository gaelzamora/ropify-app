package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	Manager  UserRole = "manager"
	attendee UserRole = "attendee"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	AvatarURL string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	Followers []*User   `json:"followers" gorm:"many2many:user_followers;joinForeignKey:UserID;JoinReferences:FollowerID"`
	Following []*User   `json:"following" gorm:"many2many:user_following;joinForeignKey:UserID;JoinReferences:FollowingID"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"-"` // No exponer el password
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
