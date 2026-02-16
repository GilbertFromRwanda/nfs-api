package agreement

import (
	"errors"
	"time"

	"nfs-api/database"
	"nfs-api/utils"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(req CreateRequest) (*Agreement, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	a := Agreement{
		Name:           req.Name,
		NumberOfCopies: req.NumberOfCopies,
		Pricing:        req.Pricing,
		Date:           date,
		Content:        req.Content,
		NotaryID:       req.NotaryID,
		OfficeID:       req.OfficeID,
	}

	if err := database.DB.Create(&a).Error; err != nil {
		return nil, errors.New("failed to create agreement")
	}
	return &a, nil
}

func (s *Service) GetAll(filter FilterRequest) (*utils.PaginatedResponse, error) {
	var agreements []Agreement
	var total int64
	query := database.DB.Model(&Agreement{})

	if filter.DateFrom != "" {
		from, err := time.Parse("2006-01-02", filter.DateFrom)
		if err != nil {
			return nil, errors.New("invalid date_from format, use YYYY-MM-DD")
		}
		query = query.Where("date >= ?", from)
	}
	if filter.DateTo != "" {
		to, err := time.Parse("2006-01-02", filter.DateTo)
		if err != nil {
			return nil, errors.New("invalid date_to format, use YYYY-MM-DD")
		}
		query = query.Where("date <= ?", to)
	}
	if filter.NotaryID != 0 {
		query = query.Where("notary_id = ?", filter.NotaryID)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.OfficeID != 0 {
		query = query.Where("office_id = ?", filter.OfficeID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count agreements")
	}

	if err := query.
		Preload("Notary", func(db *gorm.DB) *gorm.DB { return db.Select("id, names, phone") }).
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Scopes(utils.Paginate(filter.PaginationRequest)).
		Order("date DESC").
		Find(&agreements).Error; err != nil {
		return nil, errors.New("failed to fetch agreements")
	}

	result := utils.NewPaginatedResponse(agreements, total, filter.PaginationRequest)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*Agreement, error) {
	var a Agreement
	if err := database.DB.Preload("Notary", func(db *gorm.DB) *gorm.DB { return db.Select("id, names, phone") }).Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).First(&a, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("agreement not found")
		}
		return nil, errors.New("failed to fetch agreement")
	}
	return &a, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*Agreement, error) {
	a, err := s.GetByID(id)
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
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		updates["date"] = date
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.NotaryID != nil {
		updates["notary_id"] = *req.NotaryID
	}
	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}

	if len(updates) > 0 {
		if err := database.DB.Model(a).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update agreement")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&Agreement{}, id)
	if result.Error != nil {
		return errors.New("failed to delete agreement")
	}
	if result.RowsAffected == 0 {
		return errors.New("agreement not found")
	}
	return nil
}
