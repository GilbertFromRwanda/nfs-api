package agreement

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
// @Summary      Create a new agreement
// @Tags         Agreements
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=Agreement}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/agreements [post]
// @Security BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	a, err := h.service.Create(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, a)
}

// GetAll godoc
// @Summary      List all agreements
// @Tags         Agreements
// @Produce      json
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        per_page   query     int     false  "Items per page"  default(50)
// @Param        date_from  query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to    query     string  false  "Filter to date (YYYY-MM-DD)"
// @Param        notary_id  query     int     false  "Filter by notary ID"
// @Param        name       query     string  false  "Filter by name"
// @Param        office_id  query     int     false  "Filter by office ID"
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/agreements [get]
// @Security BearerAuth
func (h *Handler) GetAll(c *gin.Context) {
	var filter FilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	agreements, err := h.service.GetAll(filter)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, agreements)
}

// GetByID godoc
// @Summary      Get agreement by ID
// @Tags         Agreements
// @Produce      json
// @Param        id   path      int  true  "Agreement ID"
// @Success      200  {object}  utils.APIResponse{data=Agreement}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/agreements/{id} [get]
// @Security BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	a, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, a)
}

// Update godoc
// @Summary      Update an agreement
// @Tags         Agreements
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Agreement ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=Agreement}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/agreements/{id} [put]
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

	a, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, a)
}

// Delete godoc
// @Summary      Delete an agreement
// @Tags         Agreements
// @Produce      json
// @Param        id   path      int  true  "Agreement ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/agreements/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "agreement deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	agreements := r.Group("/agreements")
	{
		agreements.POST("", h.Create)
		agreements.GET("", h.GetAll)
		agreements.GET("/:id", h.GetByID)
		agreements.PUT("/:id", h.Update)
		agreements.DELETE("/:id", h.Delete)

		// Agreement clients (multi-party fingerprint linkage)
		agreements.GET("/:id/clients", h.ListClients)
		agreements.POST("/:id/clients", h.AddClient)
		agreements.DELETE("/:id/clients/:linkID", h.RemoveClient)
	}
}
