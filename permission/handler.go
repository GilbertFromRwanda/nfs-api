package permission

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
// @Summary      Create a new permission
// @Tags         Permissions
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=Permission}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/permissions [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	perm, err := h.service.Create(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, perm)
}

// GetAll godoc
// @Summary      List all permissions
// @Tags         Permissions
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        per_page  query     int  false  "Items per page"  default(50)
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/permissions [get]
func (h *Handler) GetAll(c *gin.Context) {
	var p utils.PaginationRequest
	if err := c.ShouldBindQuery(&p); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	perms, err := h.service.GetAll(p)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, perms)
}

// GetByID godoc
// @Summary      Get permission by ID
// @Tags         Permissions
// @Produce      json
// @Param        id   path      int  true  "Permission ID"
// @Success      200  {object}  utils.APIResponse{data=Permission}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/permissions/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	perm, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, perm)
}

// Update godoc
// @Summary      Update a permission
// @Tags         Permissions
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Permission ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=Permission}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/permissions/{id} [put]
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

	perm, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, perm)
}

// Delete godoc
// @Summary      Delete a permission
// @Tags         Permissions
// @Produce      json
// @Param        id   path      int  true  "Permission ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/permissions/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "permission deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	perms := r.Group("/permissions")
	{
		perms.POST("", h.Create)
		perms.GET("", h.GetAll)
		perms.GET("/:id", h.GetByID)
		perms.PUT("/:id", h.Update)
		perms.DELETE("/:id", h.Delete)
	}
}
