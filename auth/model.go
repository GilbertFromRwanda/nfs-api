package auth

import (
	"time"

	"nfs-api/office"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FirstName string         `gorm:"size:100;not null" json:"first_name"`
	LastName  string         `gorm:"size:100;not null" json:"last_name"`
	Username  string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	OfficeID  *uint          `gorm:"default:null" json:"office_id"`
	Office    *office.Office `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	Role      *string        `gorm:"size:50;default:null" json:"role"`
	Status    string         `gorm:"size:10;not null;default:'active'" json:"status"` // active | inactive
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	FullName  string         `gorm:"-" json:"full_name"`
}

func (u *User) AfterFind(tx *gorm.DB) error {
	u.FullName = u.FirstName + " " + u.LastName
	return nil
}
