package nlafile

import (
	"errors"
	"fmt"
	"time"

	"nfs-api/database"
	"nfs-api/utils"

	"gorm.io/gorm"
)

// Valid status transitions
var transitions = map[string]string{
	"pending":        "scanning",
	"scanning":       "scanned",
	"scanned":        "request_action",
	"request_action": "nla",
	"nla":            "approved", // default forward; rejected is an alternate
}

// Maps status to the corresponding timestamp column
var statusColumn = map[string]string{
	"pending":        "pending_at",
	"scanning":       "scanning_at",
	"scanned":        "scanned_at",
	"request_action": "request_action_at",
	"nla":            "nla_at",
	"rejected":       "rejected_at",
	"approved":       "approved_at",
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(req CreateRequest, userID uint) (*NlaFile, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	now := time.Now()
	nla := NlaFile{
		OfficeServiceID:  req.OfficeServiceID,
		Seller:           req.Seller,
		Buyer:            req.Buyer,
		Date:             date,
		PendingPeriodDays: req.PendingPeriodDays,
		UPI:              req.UPI,
		OfficeID:         req.OfficeID,
		UserID:           userID,
		Comment:          req.Comment,
		Status:           "pending",
		PendingAt:        &now,
	}

	if err := database.DB.Create(&nla).Error; err != nil {
		return nil, errors.New("failed to create nla file")
	}
	return &nla, nil
}

func (s *Service) GetAll(filter FilterRequest) (*utils.PaginatedResponse, error) {
	var files []NlaFile
	var total int64
	query := database.DB.Model(&NlaFile{})

	if filter.UPI != "" {
		query = query.Where("upi LIKE ?", "%"+filter.UPI+"%")
	}
	if filter.Name != "" {
		query = query.Where("(seller LIKE ? OR buyer LIKE ?)", "%"+filter.Name+"%", "%"+filter.Name+"%")
	}
	if filter.OfficeID != 0 {
		query = query.Where("office_id = ?", filter.OfficeID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
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

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count nla files")
	}

	if err := query.
		Preload("OfficeService", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, pricing") }).
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id, first_name, last_name, username") }).
		Scopes(utils.Paginate(filter.PaginationRequest)).
		Order("created_at DESC").
		Find(&files).Error; err != nil {
		return nil, errors.New("failed to fetch nla files")
	}

	result := utils.NewPaginatedResponse(files, total, filter.PaginationRequest)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*NlaFile, error) {
	var nla NlaFile
	if err := database.DB.
		Preload("OfficeService", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, pricing") }).
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id, first_name, last_name, username") }).
		First(&nla, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("nla file not found")
		}
		return nil, errors.New("failed to fetch nla file")
	}
	return &nla, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*NlaFile, error) {
	nla, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.OfficeServiceID != nil {
		updates["office_service_id"] = *req.OfficeServiceID
	}
	if req.Seller != nil {
		updates["seller"] = *req.Seller
	}
	if req.Buyer != nil {
		updates["buyer"] = *req.Buyer
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		updates["date"] = date
	}
	if req.PendingPeriodDays != nil {
		updates["pending_period_days"] = *req.PendingPeriodDays
	}
	if req.UPI != nil {
		updates["upi"] = *req.UPI
	}
	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}

	if len(updates) > 0 {
		if err := database.DB.Model(nla).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update nla file")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&NlaFile{}, id)
	if result.Error != nil {
		return errors.New("failed to delete nla file")
	}
	if result.RowsAffected == 0 {
		return errors.New("nla file not found")
	}
	return nil
}

// MoveToNext advances the NLA file to the next status in the workflow.
// pending -> scanning -> scanned -> request_action -> nla -> approved
func (s *Service) MoveToNext(id uint) (*NlaFile, error) {
	nla, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	nextStatus, ok := transitions[nla.Status]
	if !ok {
		return nil, fmt.Errorf("cannot advance from status '%s'", nla.Status)
	}

	return s.applyStatus(nla, nextStatus)
}

// Reject moves the NLA file to rejected status. Only allowed from "nla" status.
func (s *Service) Reject(id uint) (*NlaFile, error) {
	nla, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if nla.Status != "nla" {
		return nil, fmt.Errorf("can only reject from 'nla' status, current: '%s'", nla.Status)
	}

	return s.applyStatus(nla, "rejected")
}

// Approve moves the NLA file to approved status. Only allowed from "nla" status.
func (s *Service) Approve(id uint) (*NlaFile, error) {
	nla, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if nla.Status != "nla" {
		return nil, fmt.Errorf("can only approve from 'nla' status, current: '%s'", nla.Status)
	}

	return s.applyStatus(nla, "approved")
}

func (s *Service) GetStatusCounts() ([]StatusCount, error) {
	var counts []StatusCount
	err := database.DB.Model(&NlaFile{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&counts).Error
	if err != nil {
		return nil, errors.New("failed to fetch status counts")
	}
	return counts, nil
}

func (s *Service) applyStatus(nla *NlaFile, newStatus string) (*NlaFile, error) {
	now := time.Now()
	col, ok := statusColumn[newStatus]
	if !ok {
		return nil, fmt.Errorf("unknown status '%s'", newStatus)
	}

	updates := map[string]interface{}{
		"status": newStatus,
		col:      now,
	}

	if err := database.DB.Model(nla).Updates(updates).Error; err != nil {
		return nil, errors.New("failed to update nla file status")
	}

	return s.GetByID(nla.ID)
}
