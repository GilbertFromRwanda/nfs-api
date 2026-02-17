package nlafile

import (
	"time"

	"nfs-api/auth"
	"nfs-api/office"
	"nfs-api/officeservice"
)

type NlaFile struct {
	ID                uint                         `gorm:"primaryKey" json:"id"`
	OfficeServiceID   uint                         `gorm:"not null;index" json:"office_service_id"`
	OfficeService     *officeservice.OfficeService `gorm:"foreignKey:OfficeServiceID" json:"office_service,omitempty"`
	Seller            string                       `gorm:"size:150;not null" json:"seller"`
	Buyer             string                       `gorm:"size:150;not null" json:"buyer"`
	Date              time.Time                    `gorm:"not null;index" json:"date"`
	PendingPeriodDays int                          `gorm:"not null;default:3" json:"pending_period_days"`
	UPI               string                       `gorm:"size:40;not null;index" json:"upi"`
	OfficeID          uint                         `gorm:"not null;index" json:"office_id"`
	Office            *office.Office               `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	UserID            uint                         `gorm:"not null;index" json:"user_id"`
	User              *auth.User                   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comment           string                       `gorm:"type:text" json:"comment"`
	Status            string                       `gorm:"size:20;not null;default:'pending'" json:"status"`
	PendingAt         *time.Time                   `json:"pending_at"`
	ScanningAt        *time.Time                   `json:"scanning_at"`
	ScannedAt         *time.Time                   `json:"scanned_at"`
	RequestActionAt   *time.Time                   `json:"request_action_at"`
	NlaAt             *time.Time                   `json:"nla_at"`
	RejectedAt        *time.Time                   `json:"rejected_at"`
	ApprovedAt        *time.Time                   `json:"approved_at"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
}
