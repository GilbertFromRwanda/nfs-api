package nlachecklist

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
// @Summary      Create a new NLA checklist
// @Tags         NLA Checklists
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=NlaChecklist}
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Router       /api/v1/nla-checklists [post]
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

	checklist, err := h.service.Create(req, userID.(uint), officeID)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, checklist)
}

// GetAll godoc
// @Summary      List all NLA checklists
// @Tags         NLA Checklists
// @Produce      json
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        per_page   query     int     false  "Items per page"  default(50)
// @Param        date_from  query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to    query     string  false  "Filter to date (YYYY-MM-DD)"
// @Param        status     query     string  false  "Filter by status (completed/uncompleted)"
// @Param        office_id  query     int     false  "Filter by office ID"
// @Param        upi        query     string  false  "Filter by UPI"
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/nla-checklists [get]
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

	checklists, err := h.service.GetAll(filter)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, checklists)
}

// GetByID godoc
// @Summary      Get NLA checklist by ID
// @Tags         NLA Checklists
// @Produce      json
// @Param        id   path      int  true  "Checklist ID"
// @Success      200  {object}  utils.APIResponse{data=NlaChecklist}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/nla-checklists/{id} [get]
// @Security BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	checklist, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, checklist)
}

// Update godoc
// @Summary      Update an NLA checklist
// @Tags         NLA Checklists
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Checklist ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=NlaChecklist}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/nla-checklists/{id} [put]
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

	checklist, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, checklist)
}

// Delete godoc
// @Summary      Delete an NLA checklist
// @Tags         NLA Checklists
// @Produce      json
// @Param        id   path      int  true  "Checklist ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/nla-checklists/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "checklist deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	checklists := r.Group("/nla-checklists")
	{
		checklists.POST("", h.Create)
		checklists.GET("", h.GetAll)
		checklists.GET("/:id", h.GetByID)
		checklists.PUT("/:id", h.Update)
		checklists.DELETE("/:id", h.Delete)
	}
}
