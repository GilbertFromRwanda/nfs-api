package officeform

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
// @Summary      Create a new office form
// @Tags         Office Forms
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=OfficeForm}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/office-forms [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	form, err := h.service.Create(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, form)
}

// GetAll godoc
// @Summary      List all office forms
// @Tags         Office Forms
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        per_page  query     int  false  "Items per page"  default(50)
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/office-forms [get]
func (h *Handler) GetAll(c *gin.Context) {
	var p utils.PaginationRequest
	if err := c.ShouldBindQuery(&p); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	forms, err := h.service.GetAll(p)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, forms)
}

// GetByOfficeID godoc
// @Summary      Get office forms by office ID
// @Tags         Office Forms
// @Produce      json
// @Param        officeId  path      int  true  "Office ID"
// @Success      200       {object}  utils.APIResponse{data=[]OfficeForm}
// @Failure      400       {object}  utils.APIResponse
// @Router       /api/office-forms/office/{officeId} [get]
func (h *Handler) GetByOfficeID(c *gin.Context) {
	officeID, err := strconv.ParseUint(c.Param("officeId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid office id")
		return
	}

	forms, err := h.service.GetByOfficeID(uint(officeID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, forms)
}

// GetByID godoc
// @Summary      Get office form by ID
// @Tags         Office Forms
// @Produce      json
// @Param        id   path      int  true  "Office Form ID"
// @Success      200  {object}  utils.APIResponse{data=OfficeForm}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/office-forms/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	form, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, form)
}

// Update godoc
// @Summary      Update an office form
// @Tags         Office Forms
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Office Form ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=OfficeForm}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/office-forms/{id} [put]
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

	form, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, form)
}

// Delete godoc
// @Summary      Delete an office form
// @Tags         Office Forms
// @Produce      json
// @Param        id   path      int  true  "Office Form ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/office-forms/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "office form deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	forms := r.Group("/office-forms")
	{
		forms.POST("", h.Create)
		forms.GET("", h.GetAll)
		forms.GET("/:id", h.GetByID)
		forms.GET("/office/:officeId", h.GetByOfficeID)
		forms.PUT("/:id", h.Update)
		forms.DELETE("/:id", h.Delete)
	}
}
