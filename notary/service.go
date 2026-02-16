package notary

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

func (s *Service) Create(req CreateRequest) (*Notary, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	n := Notary{
		Names:    req.Names,
		Phone:    req.Phone,
		Status:   status,
		OfficeID: req.OfficeID,
	}

	if err := database.DB.Create(&n).Error; err != nil {
		return nil, errors.New("failed to create notary")
	}
	return &n, nil
}

func (s *Service) GetAll(p utils.PaginationRequest) (*utils.PaginatedResponse, error) {
	var notaries []Notary
	var total int64

	query := database.DB.Model(&Notary{}).Where("deleted = ?", 0)
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count notaries")
	}

	if err := query.
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Scopes(utils.Paginate(p)).
		Find(&notaries).Error; err != nil {
		return nil, errors.New("failed to fetch notaries")
	}

	result := utils.NewPaginatedResponse(notaries, total, p)
	return &result, nil
}

func (s *Service) GetByOfficeID(officeID uint) ([]Notary, error) {
	var notaries []Notary
	if err := database.DB.Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).Where("office_id = ? AND deleted = ?", officeID, 0).Find(&notaries).Error; err != nil {
		return nil, errors.New("failed to fetch notaries")
	}
	return notaries, nil
}

func (s *Service) GetByID(id uint) (*Notary, error) {
	var n Notary
	if err := database.DB.Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).Where("deleted = ?", 0).First(&n, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notary not found")
		}
		return nil, errors.New("failed to fetch notary")
	}
	return &n, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*Notary, error) {
	n, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Names != nil {
		updates["names"] = *req.Names
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}

	if len(updates) > 0 {
		if err := database.DB.Model(n).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update notary")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	n, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Model(n).Update("deleted", 1).Error; err != nil {
		return errors.New("failed to delete notary")
	}
	return nil
}
