package officepayment

import (
	"time"

	"nfs-api/office"

	"gorm.io/gorm"
)

type OfficePayment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OfficeID    uint           `gorm:"not null;index" json:"office_id"`
	Office      *office.Office `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	PaidAmount  float64        `gorm:"not null" json:"paid_amount"`
	PaidMonth   int            `gorm:"not null" json:"paid_month"`
	PaymentDate time.Time      `gorm:"not null;index" json:"payment_date"`
	EndDate     time.Time      `gorm:"not null;index" json:"end_date"`
	Comment     string         `gorm:"type:text" json:"comment"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
