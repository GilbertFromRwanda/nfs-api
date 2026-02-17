package nlachecklist

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

func (s *Service) Create(req CreateRequest, userID uint, officeID uint) (NlaChecklist, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return NlaChecklist{}, errors.New("invalid date format, use YYYY-MM-DD")
	}

	var dossiers []ChecklistDossier
	for _, d := range req.Dossiers {
		dossiers = append(dossiers, ChecklistDossier{
			Name:   d.Name,
			Status: d.Status,
		})
	}

	checklist := NlaChecklist{
		UPI:          req.UPI,
		Seller1Names: req.Seller1Names,
		Seller2Names: req.Seller2Names,
		Buyer1Names:  req.Buyer1Names,
		Buyer2Names:  req.Buyer2Names,
		UserID:       userID,
		OfficeID:     officeID,
		Date:         date,
		Dossiers:     dossiers,
	}

	if err := database.DB.Create(&checklist).Error; err != nil {
		return NlaChecklist{}, err
	}

	checklist.ComputeStatus()
	return checklist, nil
}

func (s *Service) GetAll(filter FilterRequest) (utils.PaginatedResponse, error) {
	var checklists []NlaChecklist
	var total int64

	query := database.DB.Model(&NlaChecklist{})
	query = s.applyFilters(query, filter)

	query.Count(&total)

	err := query.Scopes(utils.Paginate(filter.PaginationRequest)).
		Preload("Dossiers").
		Preload("User").
		Order("created_at DESC").
		Find(&checklists).Error
	if err != nil {
		return utils.PaginatedResponse{}, err
	}

	// Compute status and post-filter if needed
	if filter.Status != "" {
		var filtered []NlaChecklist
		for i := range checklists {
			checklists[i].ComputeStatus()
			if checklists[i].Status == filter.Status {
				filtered = append(filtered, checklists[i])
			}
		}
		return utils.NewPaginatedResponse(filtered, int64(len(filtered)), filter.PaginationRequest), nil
	}

	for i := range checklists {
		checklists[i].ComputeStatus()
	}

	return utils.NewPaginatedResponse(checklists, total, filter.PaginationRequest), nil
}

func (s *Service) GetByID(id uint) (NlaChecklist, error) {
	var checklist NlaChecklist
	err := database.DB.Preload("Dossiers").Preload("User").First(&checklist, id).Error
	if err != nil {
		return NlaChecklist{}, errors.New("checklist not found")
	}
	checklist.ComputeStatus()
	return checklist, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (NlaChecklist, error) {
	var checklist NlaChecklist
	if err := database.DB.First(&checklist, id).Error; err != nil {
		return NlaChecklist{}, errors.New("checklist not found")
	}

	updates := map[string]interface{}{}

	if req.UPI != nil {
		updates["upi"] = *req.UPI
	}
	if req.Seller1Names != nil {
		updates["seller1_names"] = *req.Seller1Names
	}
	if req.Seller2Names != nil {
		updates["seller2_names"] = *req.Seller2Names
	}
	if req.Buyer1Names != nil {
		updates["buyer1_names"] = *req.Buyer1Names
	}
	if req.Buyer2Names != nil {
		updates["buyer2_names"] = *req.Buyer2Names
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return NlaChecklist{}, errors.New("invalid date format, use YYYY-MM-DD")
		}
		updates["date"] = date
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&checklist).Updates(updates).Error; err != nil {
			return NlaChecklist{}, err
		}
	}

	// Replace dossiers if provided
	if req.Dossiers != nil {
		database.DB.Where("nla_checklist_id = ?", id).Delete(&ChecklistDossier{})
		var dossiers []ChecklistDossier
		for _, d := range req.Dossiers {
			dossiers = append(dossiers, ChecklistDossier{
				NlaChecklistID: id,
				Name:           d.Name,
				Status:         d.Status,
			})
		}
		if len(dossiers) > 0 {
			database.DB.Create(&dossiers)
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&NlaChecklist{}, id)
	if result.RowsAffected == 0 {
		return errors.New("checklist not found")
	}
	return result.Error
}

func (s *Service) applyFilters(query *gorm.DB, filter FilterRequest) *gorm.DB {
	if filter.OfficeID != 0 {
		query = query.Where("office_id = ?", filter.OfficeID)
	}
	if filter.UPI != "" {
		query = query.Where("upi LIKE ?", "%"+filter.UPI+"%")
	}
	if filter.DateFrom != "" {
		query = query.Where("date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("date <= ?", filter.DateTo)
	}
	return query
}
