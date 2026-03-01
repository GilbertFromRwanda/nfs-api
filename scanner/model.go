package scanner

import "time"

type ScannerClient struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:200;not null" json:"name"`
	Nid                 string    `gorm:"size:50;not null;uniqueIndex:idx_nid_company" json:"nid"`
	ServiceID           uint      `gorm:"not null;index" json:"service_id"`
	Service             string    `gorm:"size:200" json:"service"`
	FingerprintTemplate string    `gorm:"type:text;not null" json:"-"`
	CompanyID           uint      `gorm:"not null;index;uniqueIndex:idx_nid_company" json:"company_id"`
	UserID              uint      `gorm:"index" json:"user_id"`
	DeviceID            string    `gorm:"size:100" json:"device_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
