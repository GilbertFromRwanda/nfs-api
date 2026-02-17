package invoice

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nfs-api/database"
	"nfs-api/utils"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(req CreateRequest, userID uint, officeID uint) (*Invoice, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, use YYYY-MM-DD")
	}

	notaryStatus := req.NotaryStatus
	if notaryStatus == "" {
		notaryStatus = "pending"
	}

	services := make([]InvoiceService, len(req.Services))
	for i, svc := range req.Services {
		ps := svc.PaymentStatus
		if ps == "" {
			ps = "unpaid"
		}
		services[i] = InvoiceService{
			OfficeServiceID: svc.OfficeServiceID,
			Name:            svc.Name,
			Copy:            svc.Copy,
			Price:           svc.Price,
			PaymentStatus:   ps,
		}
	}

	inv := Invoice{
		UPI:          req.UPI,
		ClientName:   req.ClientName,
		NotaryID:     req.NotaryID,
		Date:         date,
		Note:         req.Note,
		RefNo:        req.RefNo,
		OfficeID:     officeID,
		UserID:       userID,
		NotaryStatus: notaryStatus,
		Services:     services,
	}

	if err := database.DB.Create(&inv).Error; err != nil {
		return nil, errors.New("failed to create invoice")
	}

	return s.GetByID(inv.ID)
}

func (s *Service) GetAll(filter FilterRequest) (*utils.PaginatedResponse, error) {
	var invoices []Invoice
	var total int64

	query := database.DB.Model(&Invoice{})
	query = s.applyFilters(query, filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count invoices")
	}

	if err := query.
		Preload("Services").
		Preload("Notary", func(db *gorm.DB) *gorm.DB { return db.Select("id, names, phone") }).
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id, first_name, last_name, username") }).
		Scopes(utils.Paginate(filter.PaginationRequest)).
		Order("created_at DESC").
		Find(&invoices).Error; err != nil {
		return nil, errors.New("failed to fetch invoices")
	}

	for i := range invoices {
		invoices[i].ComputeTotals()
	}

	result := utils.NewPaginatedResponse(invoices, total, filter.PaginationRequest)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*Invoice, error) {
	var inv Invoice
	if err := database.DB.
		Preload("Services").
		Preload("Notary", func(db *gorm.DB) *gorm.DB { return db.Select("id, names, phone") }).
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id, first_name, last_name, username") }).
		First(&inv, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice not found")
		}
		return nil, errors.New("failed to fetch invoice")
	}
	inv.ComputeTotals()
	return &inv, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*Invoice, error) {
	inv, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.UPI != nil {
		updates["upi"] = *req.UPI
	}
	if req.ClientName != nil {
		updates["client_name"] = *req.ClientName
	}
	if req.NotaryID != nil {
		updates["notary_id"] = *req.NotaryID
	}
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, use YYYY-MM-DD")
		}
		updates["date"] = date
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.RefNo != nil {
		updates["ref_no"] = *req.RefNo
	}
	if req.NotaryStatus != nil {
		updates["notary_status"] = *req.NotaryStatus
	}

	if len(updates) > 0 {
		if err := database.DB.Model(inv).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update invoice")
		}
	}

	return s.GetByID(id)
}

func (s *Service) UpdateServiceStatus(invoiceID, serviceID uint, req UpdateServiceStatusRequest) (*Invoice, error) {
	result := database.DB.Model(&InvoiceService{}).
		Where("id = ? AND invoice_id = ?", serviceID, invoiceID).
		Update("payment_status", req.PaymentStatus)

	if result.Error != nil {
		return nil, errors.New("failed to update service status")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("invoice service not found")
	}

	return s.GetByID(invoiceID)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&Invoice{}, id)
	if result.Error != nil {
		return errors.New("failed to delete invoice")
	}
	if result.RowsAffected == 0 {
		return errors.New("invoice not found")
	}
	return nil
}

func (s *Service) GetUserSummaries(filter FilterRequest) ([]UserSummary, error) {
	query := database.DB.Table("invoices").
		Select(`
			invoices.user_id,
			users.first_name,
			users.last_name,
			COALESCE(SUM(invoice_services.price), 0) as total_amount,
			COALESCE(SUM(CASE WHEN invoice_services.payment_status = 'paid' THEN invoice_services.price ELSE 0 END), 0) as paid_amount,
			COALESCE(SUM(CASE WHEN invoice_services.payment_status = 'unpaid' THEN invoice_services.price ELSE 0 END), 0) as unpaid_amount
		`).
		Joins("JOIN users ON users.id = invoices.user_id").
		Joins("LEFT JOIN invoice_services ON invoice_services.invoice_id = invoices.id")

	query = s.applyRawFilters(query, filter)

	query = query.Group("invoices.user_id, users.first_name, users.last_name")

	switch filter.UserStatus {
case "paid":
		query = query.Having("unpaid_amount = 0")
	case "unpaid":
		query = query.Having("unpaid_amount > 0")
	}

	var summaries []UserSummary
	if err := query.Scan(&summaries).Error; err != nil {
		return nil, errors.New("failed to fetch user summaries")
	}

	if len(summaries) == 0 {
		return summaries, nil
	}

	// Collect user IDs to fetch service breakdowns
	userIDs := make([]uint, len(summaries))
	for i, s := range summaries {
		userIDs[i] = s.UserID
	}

	// Query per-user, per-service aggregation
	type userServiceRow struct {
		UserID          uint    `json:"user_id"`
		OfficeServiceID uint    `json:"office_service_id"`
		Name            string  `json:"name"`
		TotalAmount     float64 `json:"total_amount"`
		PaidAmount      float64 `json:"paid_amount"`
		UnpaidAmount    float64 `json:"unpaid_amount"`
	}

	svcQuery := database.DB.Table("invoices").
		Select(`
			invoices.user_id,
			invoice_services.office_service_id,
			invoice_services.name,
			COALESCE(SUM(invoice_services.price), 0) as total_amount,
			COALESCE(SUM(CASE WHEN invoice_services.payment_status = 'paid' THEN invoice_services.price ELSE 0 END), 0) as paid_amount,
			COALESCE(SUM(CASE WHEN invoice_services.payment_status = 'unpaid' THEN invoice_services.price ELSE 0 END), 0) as unpaid_amount
		`).
		Joins("JOIN invoice_services ON invoice_services.invoice_id = invoices.id").
		Where("invoices.user_id IN ?", userIDs)

	svcQuery = s.applyRawFilters(svcQuery, filter)
	svcQuery = svcQuery.Group("invoices.user_id, invoice_services.office_service_id, invoice_services.name")

	var svcRows []userServiceRow
	if err := svcQuery.Scan(&svcRows).Error; err != nil {
		return nil, errors.New("failed to fetch user service summaries")
	}

	// Map services to users
	svcMap := make(map[uint][]UserServiceSummary)
	for _, r := range svcRows {
		svcMap[r.UserID] = append(svcMap[r.UserID], UserServiceSummary{
			OfficeServiceID: r.OfficeServiceID,
			Name:            r.Name,
			TotalAmount:     r.TotalAmount,
			PaidAmount:      r.PaidAmount,
			UnpaidAmount:    r.UnpaidAmount,
		})
	}

	for i := range summaries {
		summaries[i].Services = svcMap[summaries[i].UserID]
		if summaries[i].Services == nil {
			summaries[i].Services = []UserServiceSummary{}
		}
	}

	return summaries, nil
}

func (s *Service) applyFilters(query *gorm.DB, filter FilterRequest) *gorm.DB {
	if filter.OfficeID != 0 {
		query = query.Where("office_id = ?", filter.OfficeID)
	}
	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.NotaryID != 0 {
		query = query.Where("notary_id = ?", filter.NotaryID)
	}
	if filter.NotaryStatus != "" {
		query = query.Where("notary_status = ?", filter.NotaryStatus)
	}
	if filter.DateFrom != "" {
		if from, err := time.Parse("2006-01-02", filter.DateFrom); err == nil {
			query = query.Where("date >= ?", from)
		}
	}
	if filter.DateTo != "" {
		if to, err := time.Parse("2006-01-02", filter.DateTo); err == nil {
			query = query.Where("date <= ?", to)
		}
	}
	if filter.ServiceIDs != "" {
		ids := parseUintSlice(filter.ServiceIDs)
		if len(ids) > 0 {
			query = query.Where("id IN (SELECT invoice_id FROM invoice_services WHERE office_service_id IN ?)", ids)
		}
	}
	if filter.UserStatus != "" {
		if filter.UserStatus == "paid" {
			query = query.Where(`id NOT IN (
				SELECT invoice_id FROM invoice_services WHERE payment_status = 'unpaid'
			)`)
		} else if filter.UserStatus == "unpaid" {
			query = query.Where(`id IN (
				SELECT invoice_id FROM invoice_services WHERE payment_status = 'unpaid'
			)`)
		}
	}
	return query
}

func (s *Service) applyRawFilters(query *gorm.DB, filter FilterRequest) *gorm.DB {
	if filter.OfficeID != 0 {
		query = query.Where("invoices.office_id = ?", filter.OfficeID)
	}
	if filter.UserID != 0 {
		query = query.Where("invoices.user_id = ?", filter.UserID)
	}
	if filter.NotaryID != 0 {
		query = query.Where("invoices.notary_id = ?", filter.NotaryID)
	}
	if filter.NotaryStatus != "" {
		query = query.Where("invoices.notary_status = ?", filter.NotaryStatus)
	}
	if filter.DateFrom != "" {
		if from, err := time.Parse("2006-01-02", filter.DateFrom); err == nil {
			query = query.Where("invoices.date >= ?", from)
		}
	}
	if filter.DateTo != "" {
		if to, err := time.Parse("2006-01-02", filter.DateTo); err == nil {
			query = query.Where("invoices.date <= ?", to)
		}
	}
	if filter.ServiceIDs != "" {
		ids := parseUintSlice(filter.ServiceIDs)
		if len(ids) > 0 {
			query = query.Where("invoices.id IN (SELECT invoice_id FROM invoice_services WHERE office_service_id IN ?)", ids)
		}
	}
	return query
}

func parseUintSlice(s string) []uint {
	parts := strings.Split(s, ",")
	var result []uint
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if v, err := strconv.ParseUint(p, 10, 32); err == nil {
			result = append(result, uint(v))
		}
	}
	return result
}

func init() {
	// Suppress unused import warning
	_ = fmt.Sprintf
}
