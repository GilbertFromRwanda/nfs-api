package officeservice

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
// @Summary      Create a new office service
// @Tags         Office Services
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=OfficeService}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/office-services [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	os, err := h.service.Create(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, os)
}

// GetAll godoc
// @Summary      List all office services
// @Tags         Office Services
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        per_page  query     int  false  "Items per page"  default(50)
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/office-services [get]
func (h *Handler) GetAll(c *gin.Context) {
	var p utils.PaginationRequest
	if err := c.ShouldBindQuery(&p); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	services, err := h.service.GetAll(p)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, services)
}

// GetByOfficeID godoc
// @Summary      Get office services by office ID
// @Tags         Office Services
// @Produce      json
// @Param        officeId  path      int  true  "Office ID"
// @Success      200       {object}  utils.APIResponse{data=[]OfficeService}
// @Failure      400       {object}  utils.APIResponse
// @Router       /api/v1/office-services/office/{officeId} [get]
func (h *Handler) GetByOfficeID(c *gin.Context) {
	officeID, err := strconv.ParseUint(c.Param("officeId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid office id")
		return
	}

	services, err := h.service.GetByOfficeID(uint(officeID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, services)
}

// GetByID godoc
// @Summary      Get office service by ID
// @Tags         Office Services
// @Produce      json
// @Param        id   path      int  true  "Office Service ID"
// @Success      200  {object}  utils.APIResponse{data=OfficeService}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/office-services/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	os, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, os)
}

// Update godoc
// @Summary      Update an office service
// @Tags         Office Services
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Office Service ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=OfficeService}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/office-services/{id} [put]
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

	os, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, os)
}

// Delete godoc
// @Summary      Delete an office service
// @Tags         Office Services
// @Produce      json
// @Param        id   path      int  true  "Office Service ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/office-services/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "office service deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	svc := r.Group("/office-services")
	{
		svc.POST("", h.Create)
		svc.GET("", h.GetAll)
		svc.GET("/:id", h.GetByID)
		svc.GET("/office/:officeId", h.GetByOfficeID)
		svc.PUT("/:id", h.Update)
		svc.DELETE("/:id", h.Delete)
	}
}
