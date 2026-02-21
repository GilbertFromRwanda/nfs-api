package agreement

import (
	"time"

	"nfs-api/notary"
	"nfs-api/office"
)

type Agreement struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	Name           string          `gorm:"size:200;not null" json:"name"`
	NumberOfCopies int             `gorm:"not null;default:1" json:"number_of_copies"`
	Pricing        *float64        `gorm:"default:null" json:"pricing"`
	Date           time.Time       `gorm:"not null;index" json:"date"`
	Content        string          `gorm:"type:text;not null" json:"content"`
	NotaryID       uint            `gorm:"not null;index" json:"notary_id"`
	Notary         *notary.Notary  `gorm:"foreignKey:NotaryID" json:"notary,omitempty"`
	OfficeID       uint            `gorm:"not null;index" json:"office_id"`
	Office         *office.Office  `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
