package invoice

import (
	"net/http"
	"strconv"

	"nfs-api/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary      Create a new invoice
// @Tags         Invoices
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=Invoice}
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Router       /api/v1/invoices [post]
// @Security BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var officeID uint
	if oid, ok := c.Get("office_id"); ok {
		officeID = oid.(uint)
	}

	inv, err := h.service.Create(req, userID.(uint), officeID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, inv)
}

// GetAll godoc
// @Summary      List all invoices
// @Tags         Invoices
// @Produce      json
// @Param        page           query     int     false  "Page number"  default(1)
// @Param        per_page       query     int     false  "Items per page"  default(50)
// @Param        date_from      query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to        query     string  false  "Filter to date (YYYY-MM-DD)"
// @Param        user_id        query     int     false  "Filter by user ID"
// @Param        service_ids    query     string  false  "Filter by service IDs (comma-separated)"
// @Param        notary_id      query     int     false  "Filter by notary ID"
// @Param        user_status    query     string  false  "Filter by user status (paid/unpaid)"
// @Param        notary_status  query     string  false  "Filter by notary status (pending/signed_unpaid/signed_paid)"
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/invoices [get]
// @Security BearerAuth
func (h *Handler) GetAll(c *gin.Context) {
	var filter FilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if oid, ok := c.Get("office_id"); ok {
		filter.OfficeID = oid.(uint)
	}

	invoices, err := h.service.GetAll(filter)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, invoices)
}

// GetByID godoc
// @Summary      Get invoice by ID
// @Tags         Invoices
// @Produce      json
// @Param        id   path      int  true  "Invoice ID"
// @Success      200  {object}  utils.APIResponse{data=Invoice}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/invoices/{id} [get]
// @Security BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	inv, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, inv)
}

// Update godoc
// @Summary      Update an invoice
// @Tags         Invoices
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Invoice ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=Invoice}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/invoices/{id} [put]
// @Security BearerAuth
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	inv, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, inv)
}

// Delete godoc
// @Summary      Delete an invoice
// @Tags         Invoices
// @Produce      json
// @Param        id   path      int  true  "Invoice ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/invoices/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "invoice deleted"})
}

// UpdateServiceStatus godoc
// @Summary      Update invoice service payment status
// @Tags         Invoices
// @Accept       json
// @Produce      json
// @Param        id         path      int                        true  "Invoice ID"
// @Param        serviceId  path      int                        true  "Invoice Service ID"
// @Param        body       body      UpdateServiceStatusRequest  true  "Status payload"
// @Success      200  {object}  utils.APIResponse{data=Invoice}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/invoices/{id}/services/{serviceId}/status [patch]
// @Security BearerAuth
func (h *Handler) UpdateServiceStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	serviceID, err := strconv.ParseUint(c.Param("serviceId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid service id")
		return
	}

	var req UpdateServiceStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	inv, err := h.service.UpdateServiceStatus(uint(id), uint(serviceID), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, inv)
}

// GetUserSummaries godoc
// @Summary      Get invoice summaries grouped by user
// @Tags         Invoices
// @Produce      json
// @Param        page           query     int     false  "Page number"  default(1)
// @Param        per_page       query     int     false  "Items per page"  default(50)
// @Param        date_from      query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to        query     string  false  "Filter to date (YYYY-MM-DD)"
// @Param        notary_id      query     int     false  "Filter by notary ID"
// @Param        notary_status  query     string  false  "Filter by notary status"
// @Param        user_status    query     string  false  "Filter by user status (paid/unpaid)"
// @Success      200  {object}  utils.APIResponse{data=[]UserSummary}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/invoices/user-summaries [get]
// @Security BearerAuth
func (h *Handler) GetUserSummaries(c *gin.Context) {
	var filter FilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if oid, ok := c.Get("office_id"); ok {
		filter.OfficeID = oid.(uint)
	}

	summaries, err := h.service.GetUserSummaries(filter)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, summaries)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	inv := r.Group("/invoices")
	{
		inv.POST("", h.Create)
		inv.GET("", h.GetAll)
		inv.GET("/user-summaries", h.GetUserSummaries)
		inv.GET("/:id", h.GetByID)
		inv.PUT("/:id", h.Update)
		inv.DELETE("/:id", h.Delete)
		inv.PATCH("/:id/services/:serviceId/status", h.UpdateServiceStatus)
	}
}
