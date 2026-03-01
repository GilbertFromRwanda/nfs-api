package dashboard

import (
	"time"

	"nfs-api/database"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// --- Filter ---

type Filter struct {
	From     *time.Time
	To       *time.Time
	OfficeID *uint // super_admin only
}

// --- Shared types ---

type CountStats struct {
	Total     int64 `json:"total"`
	InRange   int64 `json:"in_range"`
	Today     int64 `json:"today"`
	ThisWeek  int64 `json:"this_week"`
	ThisMonth int64 `json:"this_month"`
}

type RevenueStats struct {
	Total     float64 `json:"total"`
	InRange   float64 `json:"in_range"`
	Paid      float64 `json:"paid"`
	Unpaid    float64 `json:"unpaid"`
	ThisMonth float64 `json:"this_month"`
}

type NlaStatusBreakdown struct {
	Pending       int64 `json:"pending"`
	Scanning      int64 `json:"scanning"`
	Scanned       int64 `json:"scanned"`
	RequestAction int64 `json:"request_action"`
	Nla           int64 `json:"nla"`
	Approved      int64 `json:"approved"`
	Rejected      int64 `json:"rejected"`
}

type InvoiceStatusBreakdown struct {
	Paid   int64 `json:"paid"`   // all services paid
	Unpaid int64 `json:"unpaid"` // has unpaid services
}

type ChecklistBreakdown struct {
	Completed   int64 `json:"completed"`   // all dossiers are yes/no_needed
	Uncompleted int64 `json:"uncompleted"` // at least one dossier is "no"
	Total       int64 `json:"total"`
}

type RecentItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type StuckNlaFile struct {
	ID         uint      `json:"id"`
	UPI        string    `json:"upi"`
	Status     string    `json:"status"`
	StuckSince time.Time `json:"stuck_since"` // when it entered the current status
	DaysStuck  int       `json:"days_stuck"`
}

// --- Super admin stats ---

type SuperAdminStats struct {
	FilteredOfficeID *uint                  `json:"filtered_office_id,omitempty"`
	Offices          CountStats             `json:"offices"`
	Users            CountStats             `json:"users"`
	Agreements       CountStats             `json:"agreements"`
	NlaFiles         CountStats             `json:"nla_files"`
	NlaStatus        NlaStatusBreakdown     `json:"nla_status"`      // pie chart
	Checklists       CountStats             `json:"checklists"`
	ChecklistStatus  ChecklistBreakdown     `json:"checklist_status"` // pie chart
	ChecklistNotes   CountStats             `json:"checklist_notes"`
	Invoices         CountStats             `json:"invoices"`
	InvoiceStatus    InvoiceStatusBreakdown `json:"invoice_status"`  // pie chart
	Revenue          RevenueStats           `json:"revenue"`
	Services         CountStats             `json:"services"`
	ScannerClients   CountStats             `json:"scanner_clients"`
	RecentOffices    []RecentItem           `json:"recent_offices"`
	RecentAgreements []RecentItem           `json:"recent_agreements"`
	StuckNlaFiles    []StuckNlaFile         `json:"stuck_nla_files"` // files stuck > 3 days in same status
}

// --- Office admin / employee stats ---

type OfficeStats struct {
	Staff            CountStats             `json:"staff"`
	Agreements       CountStats             `json:"agreements"`
	NlaFiles         CountStats             `json:"nla_files"`
	NlaStatus        NlaStatusBreakdown     `json:"nla_status"`      // pie chart
	Checklists       CountStats             `json:"checklists"`
	ChecklistStatus  ChecklistBreakdown     `json:"checklist_status"` // pie chart
	ChecklistNotes   CountStats             `json:"checklist_notes"`
	Invoices         CountStats             `json:"invoices"`
	InvoiceStatus    InvoiceStatusBreakdown `json:"invoice_status"`  // pie chart
	Revenue          RevenueStats           `json:"revenue"`
	Services         CountStats             `json:"services"`
	ScannerClients   CountStats             `json:"scanner_clients"`
	RecentAgreements []RecentItem           `json:"recent_agreements"`
	RecentNlaFiles   []RecentItem           `json:"recent_nla_files"`
	RecentInvoices   []RecentItem           `json:"recent_invoices"`
	RecentChecklists []RecentItem           `json:"recent_checklists"`
	StuckNlaFiles    []StuckNlaFile         `json:"stuck_nla_files"` // files stuck > 3 days in same status
}

// --- Helpers ---

func periodCounts(table, scopeWhere string, scopeArgs []interface{}, f Filter) CountStats {
	db := database.DB
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekStart := todayStart.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	base := db.Table(table)
	if scopeWhere != "" {
		base = base.Where(scopeWhere, scopeArgs...)
	}

	var c CountStats
	base.Count(&c.Total)
	base.Where("created_at >= ?", todayStart).Count(&c.Today)
	base.Where("created_at >= ?", weekStart).Count(&c.ThisWeek)
	base.Where("created_at >= ?", monthStart).Count(&c.ThisMonth)

	rangeQ := base
	if f.From != nil {
		rangeQ = rangeQ.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		rangeQ = rangeQ.Where("created_at <= ?", *f.To)
	}
	rangeQ.Count(&c.InRange)

	return c
}

func nlaStatusBreakdown(scopeWhere string, scopeArgs []interface{}, f Filter) NlaStatusBreakdown {
	db := database.DB
	base := db.Table("nla_files")
	if scopeWhere != "" {
		base = base.Where(scopeWhere, scopeArgs...)
	}
	if f.From != nil {
		base = base.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		base = base.Where("created_at <= ?", *f.To)
	}

	var b NlaStatusBreakdown
	base.Where("status = ?", "pending").Count(&b.Pending)
	base.Where("status = ?", "scanning").Count(&b.Scanning)
	base.Where("status = ?", "scanned").Count(&b.Scanned)
	base.Where("status = ?", "request_action").Count(&b.RequestAction)
	base.Where("status = ?", "nla").Count(&b.Nla)
	base.Where("status = ?", "approved").Count(&b.Approved)
	base.Where("status = ?", "rejected").Count(&b.Rejected)
	return b
}

func invoiceStatusBreakdown(scopeWhere string, scopeArgs []interface{}, f Filter) InvoiceStatusBreakdown {
	db := database.DB

	// paid = invoice has no unpaid services
	// unpaid = invoice has at least one unpaid service
	type result struct{ Count int64 }

	paidQ := db.Table("invoices i").
		Joins("LEFT JOIN invoice_services s ON s.invoice_id = i.id AND s.payment_status = 'unpaid'").
		Where("s.id IS NULL AND i.deleted_at IS NULL")
	unpaidQ := db.Table("invoices i").
		Joins("JOIN invoice_services s ON s.invoice_id = i.id AND s.payment_status = 'unpaid'").
		Where("i.deleted_at IS NULL")

	if scopeWhere != "" {
		paidQ = paidQ.Where(scopeWhere, scopeArgs...)
		unpaidQ = unpaidQ.Where(scopeWhere, scopeArgs...)
	}
	if f.From != nil {
		paidQ = paidQ.Where("i.created_at >= ?", *f.From)
		unpaidQ = unpaidQ.Where("i.created_at >= ?", *f.From)
	}
	if f.To != nil {
		paidQ = paidQ.Where("i.created_at <= ?", *f.To)
		unpaidQ = unpaidQ.Where("i.created_at <= ?", *f.To)
	}

	var b InvoiceStatusBreakdown
	paidQ.Distinct("i.id").Count(&b.Paid)
	unpaidQ.Distinct("i.id").Count(&b.Unpaid)
	return b
}

func revenueStats(scopeWhere string, scopeArgs []interface{}, f Filter) RevenueStats {
	db := database.DB
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	type sumRow struct{ Val float64 }

	base := db.Table("invoice_services s").
		Joins("JOIN invoices i ON i.id = s.invoice_id AND i.deleted_at IS NULL")
	if scopeWhere != "" {
		base = base.Where(scopeWhere, scopeArgs...)
	}

	var rev RevenueStats

	// total revenue (all time)
	var totalRow sumRow
	base.Select("COALESCE(SUM(s.price), 0) AS val").Scan(&totalRow)
	rev.Total = totalRow.Val

	// paid
	var paidRow sumRow
	base.Where("s.payment_status = 'paid'").Select("COALESCE(SUM(s.price), 0) AS val").Scan(&paidRow)
	rev.Paid = paidRow.Val

	// unpaid
	rev.Unpaid = rev.Total - rev.Paid

	// this month
	var monthRow sumRow
	base.Where("i.created_at >= ?", monthStart).Select("COALESCE(SUM(s.price), 0) AS val").Scan(&monthRow)
	rev.ThisMonth = monthRow.Val

	// in_range
	rangeQ := base
	if f.From != nil {
		rangeQ = rangeQ.Where("i.created_at >= ?", *f.From)
	}
	if f.To != nil {
		rangeQ = rangeQ.Where("i.created_at <= ?", *f.To)
	}
	var rangeRow sumRow
	rangeQ.Select("COALESCE(SUM(s.price), 0) AS val").Scan(&rangeRow)
	rev.InRange = rangeRow.Val

	return rev
}

func recentItems(table, titleCol, scopeWhere string, scopeArgs []interface{}, f Filter) []RecentItem {
	type row struct {
		ID        uint      `gorm:"column:id"`
		Title     string    `gorm:"column:title"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}
	var rows []row
	q := database.DB.Table(table).
		Select("id, "+titleCol+" AS title, created_at").
		Order("created_at DESC").
		Limit(5)
	if scopeWhere != "" {
		q = q.Where(scopeWhere, scopeArgs...)
	}
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at <= ?", *f.To)
	}
	q.Scan(&rows)

	items := make([]RecentItem, len(rows))
	for i, r := range rows {
		items[i] = RecentItem{ID: r.ID, Title: r.Title, CreatedAt: r.CreatedAt}
	}
	return items
}

func stuckNlaFiles(scopeWhere string, scopeArgs []interface{}) []StuckNlaFile {
	type row struct {
		ID          uint      `gorm:"column:id"`
		UPI         string    `gorm:"column:upi"`
		Status      string    `gorm:"column:status"`
		PendingAt   *time.Time `gorm:"column:pending_at"`
		ScanningAt  *time.Time `gorm:"column:scanning_at"`
		ScannedAt   *time.Time `gorm:"column:scanned_at"`
		RequestActionAt *time.Time `gorm:"column:request_action_at"`
		NlaAt       *time.Time `gorm:"column:nla_at"`
	}

	var rows []row
	q := database.DB.Table("nla_files").
		Select("id, upi, status, pending_at, scanning_at, scanned_at, request_action_at, nla_at").
		Where("status NOT IN ('approved','rejected')")
	if scopeWhere != "" {
		q = q.Where(scopeWhere, scopeArgs...)
	}
	q.Scan(&rows)

	cutoff := time.Now().UTC().AddDate(0, 0, -3)
	var result []StuckNlaFile
	for _, r := range rows {
		var since *time.Time
		switch r.Status {
		case "pending":
			since = r.PendingAt
		case "scanning":
			since = r.ScanningAt
		case "scanned":
			since = r.ScannedAt
		case "request_action":
			since = r.RequestActionAt
		case "nla":
			since = r.NlaAt
		}
		if since != nil && since.Before(cutoff) {
			days := int(time.Since(*since).Hours() / 24)
			result = append(result, StuckNlaFile{
				ID:         r.ID,
				UPI:        r.UPI,
				Status:     r.Status,
				StuckSince: *since,
				DaysStuck:  days,
			})
		}
	}
	return result
}

func checklistBreakdown(scopeWhere string, scopeArgs []interface{}, f Filter) ChecklistBreakdown {
	db := database.DB

	// completed = no dossier with status = 'no'
	// uncompleted = at least one dossier with status = 'no'
	completedQ := db.Table("nla_checklists c").
		Joins("LEFT JOIN checklist_dossiers d ON d.nla_checklist_id = c.id AND d.status = 'no'").
		Where("d.id IS NULL AND c.deleted_at IS NULL")
	uncompletedQ := db.Table("nla_checklists c").
		Joins("JOIN checklist_dossiers d ON d.nla_checklist_id = c.id AND d.status = 'no'").
		Where("c.deleted_at IS NULL")

	if scopeWhere != "" {
		completedQ = completedQ.Where(scopeWhere, scopeArgs...)
		uncompletedQ = uncompletedQ.Where(scopeWhere, scopeArgs...)
	}
	if f.From != nil {
		completedQ = completedQ.Where("c.created_at >= ?", *f.From)
		uncompletedQ = uncompletedQ.Where("c.created_at >= ?", *f.From)
	}
	if f.To != nil {
		completedQ = completedQ.Where("c.created_at <= ?", *f.To)
		uncompletedQ = uncompletedQ.Where("c.created_at <= ?", *f.To)
	}

	var b ChecklistBreakdown
	completedQ.Distinct("c.id").Count(&b.Completed)
	uncompletedQ.Distinct("c.id").Count(&b.Uncompleted)
	b.Total = b.Completed + b.Uncompleted
	return b
}

func officeScope(col string, officeID *uint) (string, []interface{}) {
	if officeID != nil {
		return "i.office_id = ?", []interface{}{*officeID}
	}
	return "", nil
}

// --- Service methods ---

func (s *Service) GetSuperAdminStats(f Filter) (*SuperAdminStats, error) {
	var officeScope, userScope string
	var officeArgs, userArgs []interface{}

	if f.OfficeID != nil {
		officeScope = "id = ? AND deleted_at IS NULL"
		officeArgs = []interface{}{*f.OfficeID}
		userScope = "office_id = ? AND deleted_at IS NULL"
		userArgs = []interface{}{*f.OfficeID}
	} else {
		officeScope = "deleted_at IS NULL"
		userScope = "deleted_at IS NULL"
	}

	agrScope, agrArgs := scopeByCol("office_id", f.OfficeID)
	nlaScope, nlaArgs := scopeByCol("office_id", f.OfficeID)
	invScope, invArgs := scopeByCol("office_id", f.OfficeID)
	svcScope, svcArgs := scopeByCol("office_id", f.OfficeID)
	scanScope, scanArgs := scopeByCol("company_id", f.OfficeID)

	invRevenueScope, invRevenueArgs := officeScope2(f.OfficeID)

	chkScope, chkArgs := scopeByCol("office_id", f.OfficeID)
	chkNoteScope, chkNoteArgs := scopeByCol("office_id", f.OfficeID)
	chkBreakdownScope := ""
	var chkBreakdownArgs []interface{}
	if f.OfficeID != nil {
		chkBreakdownScope = "c.office_id = ?"
		chkBreakdownArgs = []interface{}{*f.OfficeID}
	}

	stats := &SuperAdminStats{
		FilteredOfficeID: f.OfficeID,
		Offices:          periodCounts("offices", officeScope, officeArgs, f),
		Users:            periodCounts("users", userScope, userArgs, f),
		Agreements:       periodCounts("agreements", agrScope, agrArgs, f),
		NlaFiles:         periodCounts("nla_files", nlaScope, nlaArgs, f),
		NlaStatus:        nlaStatusBreakdown(nlaScope, nlaArgs, f),
		Checklists:       periodCounts("nla_checklists", chkScope, chkArgs, f),
		ChecklistStatus:  checklistBreakdown(chkBreakdownScope, chkBreakdownArgs, f),
		ChecklistNotes:   periodCounts("checklist_notes", chkNoteScope, chkNoteArgs, f),
		Invoices:         periodCounts("invoices", invScope, invArgs, f),
		InvoiceStatus:    invoiceStatusBreakdown(invScope, invArgs, f),
		Revenue:          revenueStats(invRevenueScope, invRevenueArgs, f),
		Services:         periodCounts("office_services", svcScope, svcArgs, f),
		ScannerClients:   periodCounts("scanner_clients", scanScope, scanArgs, f),
		RecentOffices:    recentItems("offices", "name", officeScope, officeArgs, f),
		RecentAgreements: recentItems("agreements", "title", agrScope, agrArgs, f),
		StuckNlaFiles:    stuckNlaFiles(nlaScope, nlaArgs),
	}
	return stats, nil
}

func (s *Service) GetOfficeStats(officeID uint, f Filter) (*OfficeStats, error) {
	scope := "office_id = ?"
	args := []interface{}{officeID}
	scanScope := "company_id = ?"
	scanArgs := []interface{}{officeID}
	invRevenueScope := "i.office_id = ?"
	invRevenueArgs := []interface{}{officeID}

	chkBreakdownScope := "c.office_id = ?"

	stats := &OfficeStats{
		Staff:            periodCounts("users", "office_id = ? AND deleted_at IS NULL", args, f),
		Agreements:       periodCounts("agreements", scope, args, f),
		NlaFiles:         periodCounts("nla_files", scope, args, f),
		NlaStatus:        nlaStatusBreakdown(scope, args, f),
		Checklists:       periodCounts("nla_checklists", scope, args, f),
		ChecklistStatus:  checklistBreakdown(chkBreakdownScope, args, f),
		ChecklistNotes:   periodCounts("checklist_notes", scope, args, f),
		Invoices:         periodCounts("invoices", scope, args, f),
		InvoiceStatus:    invoiceStatusBreakdown("i.office_id = ?", args, f),
		Revenue:          revenueStats(invRevenueScope, invRevenueArgs, f),
		Services:         periodCounts("office_services", scope, args, f),
		ScannerClients:   periodCounts("scanner_clients", scanScope, scanArgs, f),
		RecentAgreements: recentItems("agreements", "title", scope, args, f),
		RecentNlaFiles:   recentItems("nla_files", "upi", scope, args, f),
		RecentInvoices:   recentItems("invoices", "client_name", scope, args, f),
		RecentChecklists: recentItems("nla_checklists", "upi", scope, args, f),
		StuckNlaFiles:    stuckNlaFiles(scope, args),
	}
	return stats, nil
}

func scopeByCol(col string, id *uint) (string, []interface{}) {
	if id != nil {
		return col + " = ?", []interface{}{*id}
	}
	return "", nil
}

func officeScope2(officeID *uint) (string, []interface{}) {
	if officeID != nil {
		return "i.office_id = ?", []interface{}{*officeID}
	}
	return "", nil
}
