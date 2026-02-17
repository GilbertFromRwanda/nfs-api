package invoice

import (
	"time"

	"nfs-api/auth"
	"nfs-api/notary"
	"nfs-api/office"

	"gorm.io/gorm"
)

type Invoice struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	UPI          string           `gorm:"size:40;not null;index" json:"upi"`
	ClientName   string           `gorm:"size:150;not null" json:"client_name"`
	NotaryID     uint             `gorm:"not null;index" json:"notary_id"`
	Notary       *notary.Notary   `gorm:"foreignKey:NotaryID" json:"notary,omitempty"`
	Date         time.Time        `gorm:"not null;index" json:"date"`
	Note         string           `gorm:"type:text" json:"note"`
	RefNo        string           `gorm:"size:50" json:"ref_no"`
	OfficeID     uint             `gorm:"not null;index" json:"office_id"`
	Office       *office.Office   `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	UserID       uint             `gorm:"not null;index" json:"user_id"`
	User         *auth.User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	NotaryStatus string           `gorm:"size:20;not null;default:'pending'" json:"notary_status"` // pending | signed_unpaid | signed_paid
	Services     []InvoiceService `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE" json:"services"`

	// Computed at runtime
	TotalAmount  float64 `gorm:"-" json:"total_amount"`
	PaidAmount   float64 `gorm:"-" json:"paid_amount"`
	UnpaidAmount float64 `gorm:"-" json:"unpaid_amount"`
	UserStatus   string  `gorm:"-" json:"user_status"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type InvoiceService struct {
	ID              uint    `gorm:"primaryKey" json:"id"`
	InvoiceID       uint    `gorm:"not null;index" json:"invoice_id"`
	OfficeServiceID uint    `gorm:"not null" json:"office_service_id"`
	Name            string  `gorm:"size:100;not null" json:"name"`
	Copy            int     `gorm:"not null;default:1" json:"copy"`
	Price           float64 `gorm:"not null" json:"price"`
	PaymentStatus   string  `gorm:"size:10;not null;default:'unpaid'" json:"payment_status"` // paid | unpaid
}

func (inv *Invoice) ComputeTotals() {
	inv.TotalAmount = 0
	inv.PaidAmount = 0
	for _, s := range inv.Services {
		// amount := float64(s.Copy) * s.Price
		amount := s.Price
		inv.TotalAmount += amount
		if s.PaymentStatus == "paid" {
			inv.PaidAmount += amount
		}
	}
	inv.UnpaidAmount = inv.TotalAmount - inv.PaidAmount
	if inv.UnpaidAmount <= 0 {
		inv.UserStatus = "paid"
	} else {
		inv.UserStatus = "unpaid"
	}
}
