package office

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

func (s *Service) Create(req CreateRequest) (*Office, error) {
	status := req.Status
	if status == "" {
		status = "active"
	}

	office := Office{
		Name:        req.Name,
		Phone:       req.Phone,
		ScannerID:   req.ScannerID,
		PaymentType: req.PaymentType,
		Pricing:     req.Pricing,
		Location:    req.Location,
		Status:      status,
	}

	if req.SubscriptionEndDate != nil {
		t, err := time.Parse("2006-01-02", *req.SubscriptionEndDate)
		if err != nil {
			return nil, errors.New("invalid subscription_end_date format, use YYYY-MM-DD")
		}
		office.SubscriptionEndDate = &t
	}

	if err := database.DB.Create(&office).Error; err != nil {
		return nil, errors.New("failed to create office")
	}
	return &office, nil
}

func (s *Service) GetAll(p utils.PaginationRequest) (*utils.PaginatedResponse, error) {
	var offices []Office
	var total int64

	query := database.DB.Model(&Office{})
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count offices")
	}

	if err := query.Scopes(utils.Paginate(p)).Find(&offices).Error; err != nil {
		return nil, errors.New("failed to fetch offices")
	}

	result := utils.NewPaginatedResponse(offices, total, p)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*Office, error) {
	var office Office
	if err := database.DB.First(&office, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("office not found")
		}
		return nil, errors.New("failed to fetch office")
	}
	return &office, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*Office, error) {
	office, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.ScannerID != nil {
		updates["scanner_id"] = *req.ScannerID
	}
	if req.PaymentType != nil {
		updates["payment_type"] = *req.PaymentType
	}
	if req.Pricing != nil {
		updates["pricing"] = *req.Pricing
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.SubscriptionEndDate != nil {
		t, err := time.Parse("2006-01-02", *req.SubscriptionEndDate)
		if err != nil {
			return nil, errors.New("invalid subscription_end_date format, use YYYY-MM-DD")
		}
		updates["subscription_end_date"] = t
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(office).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update office")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&Office{}, id)
	if result.Error != nil {
		return errors.New("failed to delete office")
	}
	if result.RowsAffected == 0 {
		return errors.New("office not found")
	}
	return nil
}
