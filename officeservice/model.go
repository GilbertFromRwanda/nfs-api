package officeservice

import (
	"time"

	"nfs-api/office"
)

type OfficeService struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	NumberOfCopies int            `gorm:"not null;default:1" json:"number_of_copies"`
	Pricing        float64        `gorm:"default:0" json:"pricing"`
	OfficeID       uint           `gorm:"not null;index" json:"office_id"`
	Office         *office.Office `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	Deleted        int       `gorm:"default:0" json:"deleted"` // 0 = active, 1 = deleted
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
