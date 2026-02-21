package officepayment

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
// @Summary      Create a new office payment
// @Tags         Office Payments
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=OfficePayment}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/office-payments [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	payment, err := h.service.Create(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, payment)
}

// GetAll godoc
// @Summary      List all office payments
// @Tags         Office Payments
// @Produce      json
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        per_page   query     int     false  "Items per page"  default(50)
// @Param        office_id  query     int     false  "Filter by office ID"
// @Param        date_from  query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to    query     string  false  "Filter to date (YYYY-MM-DD)"
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/office-payments [get]
func (h *Handler) GetAll(c *gin.Context) {
	var filter FilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	payments, err := h.service.GetAll(filter)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, payments)
}

// GetByID godoc
// @Summary      Get office payment by ID
// @Tags         Office Payments
// @Produce      json
// @Param        id   path      int  true  "Office Payment ID"
// @Success      200  {object}  utils.APIResponse{data=OfficePayment}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/office-payments/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	payment, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, payment)
}

// Update godoc
// @Summary      Update an office payment
// @Tags         Office Payments
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Office Payment ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=OfficePayment}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/office-payments/{id} [put]
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

	payment, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, payment)
}

// Delete godoc
// @Summary      Delete an office payment
// @Tags         Office Payments
// @Produce      json
// @Param        id   path      int  true  "Office Payment ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/office-payments/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "office payment deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	payments := r.Group("/office-payments")
	{
		payments.POST("", h.Create)
		payments.GET("", h.GetAll)
		payments.GET("/:id", h.GetByID)
		payments.PUT("/:id", h.Update)
		payments.DELETE("/:id", h.Delete)
	}
}
