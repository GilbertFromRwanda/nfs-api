package agreement

import (
	"errors"

	"nfs-api/database"
	"nfs-api/scanner"

	"gorm.io/gorm"
)

// ListClients returns all clients linked to the given agreement.
func (s *Service) ListClients(agreementID uint) ([]AgreementClient, error) {
	var clients []AgreementClient
	err := database.DB.
		Where("agreement_id = ?", agreementID).
		Preload("ScannerClient", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, nid, service, company_id, created_at")
		}).
		Order("created_at ASC").
		Find(&clients).Error
	if err != nil {
		return nil, errors.New("failed to fetch agreement clients")
	}
	return clients, nil
}

// AddClient links a ScannerClient to an Agreement.
// Prevents duplicates and enforces office scoping.
func (s *Service) AddClient(agreementID uint, officeID uint, req AddClientRequest) (*AgreementClient, error) {
	// Verify agreement belongs to the caller's office
	var agr Agreement
	q := database.DB.Select("id, office_id")
	if officeID != 0 {
		q = q.Where("office_id = ?", officeID)
	}
	if err := q.First(&agr, agreementID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("agreement not found")
		}
		return nil, errors.New("failed to fetch agreement")
	}

	// Verify scanner client belongs to the same office
	var sc scanner.ScannerClient
	scq := database.DB.Select("id, company_id")
	if officeID != 0 {
		scq = scq.Where("company_id = ?", officeID)
	}
	if err := scq.First(&sc, req.ScannerClientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("client not found")
		}
		return nil, errors.New("failed to fetch scanner client")
	}

	// Prevent duplicate
	var existing AgreementClient
	err := database.DB.
		Where("agreement_id = ? AND scanner_client_id = ?", agreementID, req.ScannerClientID).
		First(&existing).Error
	if err == nil {
		return nil, errors.New("client already added to this agreement")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check existing client")
	}

	ac := AgreementClient{
		AgreementID:     agreementID,
		ScannerClientID: req.ScannerClientID,
		Role:            req.Role,
	}
	if err := database.DB.Create(&ac).Error; err != nil {
		return nil, errors.New("failed to add client to agreement")
	}

	// Return with preloaded scanner client
	database.DB.Preload("ScannerClient").First(&ac, ac.ID)
	return &ac, nil
}

// RemoveClient removes a client link from an agreement.
func (s *Service) RemoveClient(agreementID uint, officeID uint, linkID uint) error {
	// Verify agreement belongs to caller's office first
	var agr Agreement
	q := database.DB.Select("id, office_id")
	if officeID != 0 {
		q = q.Where("office_id = ?", officeID)
	}
	if err := q.First(&agr, agreementID).Error; err != nil {
		return errors.New("agreement not found")
	}

	result := database.DB.
		Where("id = ? AND agreement_id = ?", linkID, agreementID).
		Delete(&AgreementClient{})
	if result.Error != nil {
		return errors.New("failed to remove client")
	}
	if result.RowsAffected == 0 {
		return errors.New("client link not found")
	}
	return nil
}
