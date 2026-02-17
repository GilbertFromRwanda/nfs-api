package checklistnote

import (
	"errors"

	"nfs-api/database"
	"nfs-api/utils"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(req CreateRequest) (ChecklistNote, error) {
	note := ChecklistNote{
		Name:     req.Name,
		OfficeID: req.OfficeID,
	}

	if err := database.DB.Create(&note).Error; err != nil {
		return ChecklistNote{}, err
	}

	return note, nil
}

func (s *Service) GetAll(filter FilterRequest) (utils.PaginatedResponse, error) {
	var notes []ChecklistNote
	var total int64

	query := database.DB.Model(&ChecklistNote{})

	if filter.OfficeID != 0 {
		query = query.Where("office_id = ?", filter.OfficeID)
	}

	query.Count(&total)

	err := query.Scopes(utils.Paginate(filter.PaginationRequest)).
		Order("created_at DESC").
		Find(&notes).Error
	if err != nil {
		return utils.PaginatedResponse{}, err
	}

	return utils.NewPaginatedResponse(notes, total, filter.PaginationRequest), nil
}

func (s *Service) GetByID(id uint) (ChecklistNote, error) {
	var note ChecklistNote
	if err := database.DB.First(&note, id).Error; err != nil {
		return ChecklistNote{}, errors.New("checklist note not found")
	}
	return note, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (ChecklistNote, error) {
	var note ChecklistNote
	if err := database.DB.First(&note, id).Error; err != nil {
		return ChecklistNote{}, errors.New("checklist note not found")
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&note).Updates(updates).Error; err != nil {
			return ChecklistNote{}, err
		}
	}

	database.DB.First(&note, id)
	return note, nil
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&ChecklistNote{}, id)
	if result.RowsAffected == 0 {
		return errors.New("checklist note not found")
	}
	return result.Error
}
