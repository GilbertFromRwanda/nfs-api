package agreement

import (
	"time"

	"nfs-api/scanner"
)

// AgreementClient links a ScannerClient to an Agreement with an optional role.
// role examples: "seller", "buyer", "guarantor", "witness", etc.
type AgreementClient struct {
	ID              uint                   `gorm:"primaryKey" json:"id"`
	AgreementID     uint                   `gorm:"not null;index" json:"agreement_id"`
	ScannerClientID uint                   `gorm:"not null;index" json:"scanner_client_id"`
	Role            string                 `gorm:"size:100" json:"role"`
	ScannerClient   *scanner.ScannerClient `gorm:"foreignKey:ScannerClientID" json:"client,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}
