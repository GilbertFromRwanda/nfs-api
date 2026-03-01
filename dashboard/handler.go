package dashboard

import (
	"net/http"
	"time"

	"nfs-api/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type queryParams struct {
	From     string `form:"from"`      // YYYY-MM-DD
	To       string `form:"to"`        // YYYY-MM-DD
	OfficeID *uint  `form:"office_id"` // super_admin only
}

func parseFilter(c *gin.Context, isSuperAdmin bool) Filter {
	var q queryParams
	c.ShouldBindQuery(&q)

	var f Filter

	if q.From != "" {
		if t, err := time.Parse("2006-01-02", q.From); err == nil {
			f.From = &t
		}
	}
	if q.To != "" {
		if t, err := time.Parse("2006-01-02", q.To); err == nil {
			// end of day
			eod := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			f.To = &eod
		}
	}
	// only super_admin can filter by office_id
	if isSuperAdmin && q.OfficeID != nil {
		f.OfficeID = q.OfficeID
	}

	return f
}

// Get godoc
// @Summary      Get dashboard stats (role-based, filterable)
// @Tags         Dashboard
// @Produce      json
// @Param        from       query     string  false  "From date (YYYY-MM-DD)"
// @Param        to         query     string  false  "To date (YYYY-MM-DD)"
// @Param        office_id  query     int     false  "Filter by office (super_admin only)"
// @Success      200  {object}  utils.APIResponse
// @Failure      401  {object}  utils.APIResponse
// @Failure      403  {object}  utils.APIResponse
// @Router       /api/v1/dashboard [get]
// @Security BearerAuth
func (h *Handler) Get(c *gin.Context) {
	role, _ := c.Get("role")
	isSuperAdmin := role == "super_admin"

	f := parseFilter(c, isSuperAdmin)

	if isSuperAdmin {
		stats, err := h.service.GetSuperAdminStats(f)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to fetch stats")
			return
		}
		utils.Success(c, http.StatusOK, stats)
		return
	}

	officeID, exists := c.Get("office_id")
	if !exists {
		utils.Error(c, http.StatusForbidden, "no office associated with your account")
		return
	}

	stats, err := h.service.GetOfficeStats(officeID.(uint), f)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to fetch stats")
		return
	}
	utils.Success(c, http.StatusOK, stats)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/dashboard", h.Get)
}
