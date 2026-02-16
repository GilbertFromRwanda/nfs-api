package notary

import (
	"time"

	"nfs-api/office"
)

type Notary struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Names     string         `gorm:"size:200;not null" json:"names"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Status    string         `gorm:"size:10;not null;default:'active'" json:"status"` // active | inactive
	Deleted   int            `gorm:"default:0" json:"deleted"`                        // 0 = active, 1 = deleted
	OfficeID  uint           `gorm:"not null;index" json:"office_id"`
	Office    *office.Office `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
