package nlafile

import (
	"time"

	"nfs-api/office"
	"nfs-api/officeservice"
)

type NlaFile struct {
	ID               uint                         `gorm:"primaryKey" json:"id"`
	OfficeServiceID  uint                         `gorm:"not null;index" json:"office_service_id"`
	OfficeService    *officeservice.OfficeService  `gorm:"foreignKey:OfficeServiceID" json:"office_service,omitempty"`
	Seller           string                       `gorm:"size:200;not null" json:"seller"`
	Buyer            string                       `gorm:"size:200;not null" json:"buyer"`
	Date             time.Time                    `gorm:"not null;index" json:"date"`
	PendingPeriodDays int                         `gorm:"not null;default:3" json:"pending_period_days"`
	UPI              string                       `gorm:"size:100;not null;index" json:"upi"`
	OfficeID         uint                         `gorm:"not null;index" json:"office_id"`
	Office           *office.Office               `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	Status           string                       `gorm:"size:20;not null;default:'pending'" json:"status"`
	PendingAt        *time.Time                   `json:"pending_at"`
	ScanningAt       *time.Time                   `json:"scanning_at"`
	ScannedAt        *time.Time                   `json:"scanned_at"`
	RequestActionAt  *time.Time                   `json:"request_action_at"`
	NlaAt            *time.Time                   `json:"nla_at"`
	RejectedAt       *time.Time                   `json:"rejected_at"`
	ApprovedAt       *time.Time                   `json:"approved_at"`
	CreatedAt        time.Time                    `json:"created_at"`
	UpdatedAt        time.Time                    `json:"updated_at"`
}
