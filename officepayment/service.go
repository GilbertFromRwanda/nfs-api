package officepayment

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

func (s *Service) Create(req CreateRequest) (*OfficePayment, error) {
	paymentDate, err := time.Parse("2006-01-02", req.PaymentDate)
	if err != nil {
		return nil, errors.New("invalid payment_date format, use YYYY-MM-DD")
	}

	endDate := paymentDate.AddDate(0, req.PaidMonth, 0)

	payment := OfficePayment{
		OfficeID:    req.OfficeID,
		PaidAmount:  req.PaidAmount,
		PaidMonth:   req.PaidMonth,
		PaymentDate: paymentDate,
		EndDate:     endDate,
		Comment:     req.Comment,
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		return nil, errors.New("failed to create office payment")
	}

	return s.GetByID(payment.ID)
}

func (s *Service) GetAll(filter FilterRequest) (*utils.PaginatedResponse, error) {
	var payments []OfficePayment
	var total int64

	query := database.DB.Model(&OfficePayment{})

	if filter.OfficeID != 0 {
		query = query.Where("office_id = ?", filter.OfficeID)
	}
	if filter.DateFrom != "" {
		if from, err := time.Parse("2006-01-02", filter.DateFrom); err == nil {
			query = query.Where("payment_date >= ?", from)
		}
	}
	if filter.DateTo != "" {
		if to, err := time.Parse("2006-01-02", filter.DateTo); err == nil {
			query = query.Where("payment_date <= ?", to)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count office payments")
	}

	if err := query.
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Scopes(utils.Paginate(filter.PaginationRequest)).
		Order("payment_date DESC").
		Find(&payments).Error; err != nil {
		return nil, errors.New("failed to fetch office payments")
	}

	result := utils.NewPaginatedResponse(payments, total, filter.PaginationRequest)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*OfficePayment, error) {
	var payment OfficePayment
	if err := database.DB.
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		First(&payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("office payment not found")
		}
		return nil, errors.New("failed to fetch office payment")
	}
	return &payment, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*OfficePayment, error) {
	payment, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	needRecalc := false

	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}
	if req.PaidAmount != nil {
		updates["paid_amount"] = *req.PaidAmount
	}
	if req.PaidMonth != nil {
		updates["paid_month"] = *req.PaidMonth
		needRecalc = true
	}
	if req.PaymentDate != nil {
		date, err := time.Parse("2006-01-02", *req.PaymentDate)
		if err != nil {
			return nil, errors.New("invalid payment_date format, use YYYY-MM-DD")
		}
		updates["payment_date"] = date
		needRecalc = true
	}
	if req.Comment != nil {
		updates["comment"] = *req.Comment
	}

	// Recalculate end_date if payment_date or paid_month changed
	if needRecalc {
		paymentDate := payment.PaymentDate
		paidMonth := payment.PaidMonth

		if d, ok := updates["payment_date"]; ok {
			paymentDate = d.(time.Time)
		}
		if m, ok := updates["paid_month"]; ok {
			paidMonth = *req.PaidMonth
			_ = m
		}
		updates["end_date"] = paymentDate.AddDate(0, paidMonth, 0)
	}

	if len(updates) > 0 {
		if err := database.DB.Model(payment).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update office payment")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&OfficePayment{}, id)
	if result.Error != nil {
		return errors.New("failed to delete office payment")
	}
	if result.RowsAffected == 0 {
		return errors.New("office payment not found")
	}
	return nil
}
