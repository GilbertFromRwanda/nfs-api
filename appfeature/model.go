package appfeature

import "time"

type AppFeature struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Status    string    `gorm:"size:10;not null;default:'active'" json:"status"` // active | inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
