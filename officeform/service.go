package officeform

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

func (s *Service) Create(req CreateRequest) (*OfficeForm, error) {
	form := OfficeForm{
		Title:    req.Title,
		Content:  req.Content,
		OfficeID: req.OfficeID,
	}

	if err := database.DB.Create(&form).Error; err != nil {
		return nil, errors.New("failed to create office form")
	}
	return &form, nil
}

func (s *Service) GetAll(p utils.PaginationRequest) (*utils.PaginatedResponse, error) {
	var forms []OfficeForm
	var total int64

	query := database.DB.Model(&OfficeForm{})
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count office forms")
	}

	if err := query.Scopes(utils.Paginate(p)).Find(&forms).Error; err != nil {
		return nil, errors.New("failed to fetch office forms")
	}

	result := utils.NewPaginatedResponse(forms, total, p)
	return &result, nil
}

func (s *Service) GetByOfficeID(officeID uint) ([]OfficeForm, error) {
	var forms []OfficeForm
	if err := database.DB.Where("office_id = ?", officeID).Find(&forms).Error; err != nil {
		return nil, errors.New("failed to fetch office forms")
	}
	return forms, nil
}

func (s *Service) GetByID(id uint) (*OfficeForm, error) {
	var form OfficeForm
	if err := database.DB.First(&form, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("office form not found")
		}
		return nil, errors.New("failed to fetch office form")
	}
	return &form, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*OfficeForm, error) {
	form, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}

	if len(updates) > 0 {
		if err := database.DB.Model(form).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update office form")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&OfficeForm{}, id)
	if result.Error != nil {
		return errors.New("failed to delete office form")
	}
	if result.RowsAffected == 0 {
		return errors.New("office form not found")
	}
	return nil
}
