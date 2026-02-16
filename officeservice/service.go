package officeservice

import (
	"errors"

	"nfs-api/database"
	"nfs-api/utils"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(req CreateRequest) (*OfficeService, error) {
	os := OfficeService{
		Name:           req.Name,
		NumberOfCopies: req.NumberOfCopies,
		Pricing:        req.Pricing,
		OfficeID:       req.OfficeID,
	}

	if err := database.DB.Create(&os).Error; err != nil {
		return nil, errors.New("failed to create office service")
	}
	return &os, nil
}

func (s *Service) GetAll(p utils.PaginationRequest) (*utils.PaginatedResponse, error) {
	var services []OfficeService
	var total int64

	query := database.DB.Model(&OfficeService{}).Where("deleted = ?", 0)
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count office services")
	}

	if err := query.
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Scopes(utils.Paginate(p)).
		Find(&services).Error; err != nil {
		return nil, errors.New("failed to fetch office services")
	}

	result := utils.NewPaginatedResponse(services, total, p)
	return &result, nil
}

func (s *Service) GetByOfficeID(officeID uint) ([]OfficeService, error) {
	var services []OfficeService
	if err := database.DB.Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).Where("office_id = ? AND deleted = ?", officeID, 0).Find(&services).Error; err != nil {
		return nil, errors.New("failed to fetch office services")
	}
	return services, nil
}

func (s *Service) GetByID(id uint) (*OfficeService, error) {
	var os OfficeService
	if err := database.DB.Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).Where("deleted = ?", 0).First(&os, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("office service not found")
		}
		return nil, errors.New("failed to fetch office service")
	}
	return &os, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*OfficeService, error) {
	os, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.NumberOfCopies != nil {
		updates["number_of_copies"] = *req.NumberOfCopies
	}
	if req.Pricing != nil {
		updates["pricing"] = *req.Pricing
	}
	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}

	if len(updates) > 0 {
		if err := database.DB.Model(os).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update office service")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	os, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Model(os).Update("deleted", 1).Error; err != nil {
		return errors.New("failed to delete office service")
	}
	return nil
}
