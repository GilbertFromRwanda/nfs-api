package office

import "time"

type Office struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:100;not null" json:"name"`
	Phone               string    `gorm:"size:20" json:"phone"`
	ScannerID           string    `gorm:"size:100" json:"scanner_id"`
	SubscriptionEndDate *time.Time `json:"subscription_end_date"`
	PaymentType         string    `gorm:"size:20;not null;default:'subscription'" json:"payment_type"` // subscription | full_paid
	Pricing             float64   `gorm:"default:0" json:"pricing"`
	Location            string    `gorm:"size:255" json:"location"`
	Status              string    `gorm:"size:10;not null;default:'active'" json:"status"` // active | inactive
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
