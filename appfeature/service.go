package appfeature

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

func (s *Service) Create(req CreateRequest) (*AppFeature, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	feature := AppFeature{
		Name:   req.Name,
		Status: status,
	}

	if err := database.DB.Create(&feature).Error; err != nil {
		return nil, errors.New("failed to create app feature")
	}
	return &feature, nil
}

func (s *Service) GetAll(p utils.PaginationRequest) (*utils.PaginatedResponse, error) {
	var features []AppFeature
	var total int64

	query := database.DB.Model(&AppFeature{})
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count app features")
	}

	if err := query.Scopes(utils.Paginate(p)).Order("name ASC").Find(&features).Error; err != nil {
		return nil, errors.New("failed to fetch app features")
	}

	result := utils.NewPaginatedResponse(features, total, p)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*AppFeature, error) {
	var feature AppFeature
	if err := database.DB.First(&feature, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("app feature not found")
		}
		return nil, errors.New("failed to fetch app feature")
	}
	return &feature, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*AppFeature, error) {
	feature, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(feature).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update app feature")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&AppFeature{}, id)
	if result.Error != nil {
		return errors.New("failed to delete app feature")
	}
	if result.RowsAffected == 0 {
		return errors.New("app feature not found")
	}
	return nil
}
