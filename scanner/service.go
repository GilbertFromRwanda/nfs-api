package scanner

import (
	"errors"
	"math"

	"nfs-api/database"
	"nfs-api/utils"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Enroll(req EnrollRequest) (*ScannerClient, error) {
	client := ScannerClient{
		Name:                req.Name,
		Nid:                 req.Nid,
		Phone:               req.Phone,
		ServiceID:           req.ServiceID,
		Service:             req.Service,
		FingerprintTemplate: req.FingerprintTemplate,
		CompanyID:           req.CompanyID,
		UserID:              req.UserID,
		DeviceID:            req.DeviceID,
	}

	if err := database.DB.Create(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("this person is already enrolled in your office")
		}
		return nil, errors.New("failed to enroll client")
	}
	return &client, nil
}

func (s *Service) List(req ListRequest) (*utils.PaginatedResponse, error) {
	var clients []ScannerClient
	var total int64

	query := database.DB.Model(&ScannerClient{})

	if req.CompanyID != 0 {
		query = query.Where("company_id = ?", req.CompanyID)
	}
	if req.Q != "" {
		like := "%" + req.Q + "%"
		query = query.Where("name ILIKE ? OR nid ILIKE ?", like, like)
	}
	if req.Service != "" {
		query = query.Where("service ILIKE ?", "%"+req.Service+"%")
	}
	if req.From != "" {
		query = query.Where("created_at >= ?", req.From)
	}
	if req.To != "" {
		query = query.Where("created_at <= ?", req.To)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count scanner clients")
	}

	p := utils.PaginationRequest{Page: req.Page, PerPage: req.PerPage}
	if err := query.Scopes(utils.Paginate(p)).Find(&clients).Error; err != nil {
		return nil, errors.New("failed to fetch scanner clients")
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PerPage)))
	result := utils.PaginatedResponse{
		Items:      clients,
		Total:      total,
		Page:       req.Page,
		PerPage:    req.PerPage,
		TotalPages: totalPages,
	}
	return &result, nil
}

// UpdateFingerprint replaces the fingerprint template for an existing enrolled client.
func (s *Service) UpdateFingerprint(id uint, officeID uint, req UpdateFingerprintRequest) error {
	updates := map[string]interface{}{
		"fingerprint_template": req.FingerprintTemplate,
	}
	if req.DeviceID != "" {
		updates["device_id"] = req.DeviceID
	}

	query := database.DB.Model(&ScannerClient{})
	if officeID != 0 {
		query = query.Where("company_id = ?", officeID)
	}
	result := query.Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return errors.New("failed to update fingerprint")
	}
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}

// GetFingerprint fetches the template and enforces that the client belongs to the caller's office.
func (s *Service) GetFingerprint(id uint, officeID uint) (string, error) {
	var client ScannerClient
	query := database.DB.Select("id, fingerprint_template, company_id")
	if officeID != 0 {
		query = query.Where("company_id = ?", officeID)
	}
	if err := query.First(&client, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("not found")
		}
		return "", errors.New("internal")
	}
	return client.FingerprintTemplate, nil
}

// Get fetches a single client by ID scoped to the caller's office.
func (s *Service) Get(id uint, officeID uint) (*ScannerClient, error) {
	var client ScannerClient
	query := database.DB.Model(&ScannerClient{})
	if officeID != 0 {
		query = query.Where("company_id = ?", officeID)
	}
	if err := query.First(&client, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not found")
		}
		return nil, errors.New("internal")
	}
	return &client, nil
}

// AllTemplates returns id+name+nid+template for every enrolled client in the office.
// Used by the frontend to detect duplicate fingerprints before enrolling a new client.
func (s *Service) AllTemplates(officeID uint) ([]TemplateEntry, error) {
	var rows []TemplateEntry
	err := database.DB.Model(&ScannerClient{}).
		Select("id, name, nid, fingerprint_template AS template").
		Where("company_id = ?", officeID).
		Find(&rows).Error
	if err != nil {
		return nil, errors.New("failed to fetch templates")
	}
	return rows, nil
}

// Delete removes a client scoped to the caller's office.
func (s *Service) Delete(id uint, officeID uint) error {
	query := database.DB.Model(&ScannerClient{})
	if officeID != 0 {
		query = query.Where("company_id = ?", officeID)
	}
	result := query.Delete(&ScannerClient{}, id)
	if result.Error != nil {
		return errors.New("failed to delete client")
	}
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}
	return nil
}
