package notary

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
// @Summary      Create a new notary
// @Tags         Notaries
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=Notary}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/notaries [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	n, err := h.service.Create(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, n)
}

// GetAll godoc
// @Summary      List all notaries
// @Tags         Notaries
// @Produce      json
// @Param        page      query     int  false  "Page number"  default(1)
// @Param        per_page  query     int  false  "Items per page"  default(50)
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/notaries [get]
func (h *Handler) GetAll(c *gin.Context) {
	var p utils.PaginationRequest
	if err := c.ShouldBindQuery(&p); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	notaries, err := h.service.GetAll(p)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, notaries)
}

// GetByOfficeID godoc
// @Summary      Get notaries by office ID
// @Tags         Notaries
// @Produce      json
// @Param        officeId  path      int  true  "Office ID"
// @Success      200       {object}  utils.APIResponse{data=[]Notary}
// @Failure      400       {object}  utils.APIResponse
// @Router       /api/notaries/office/{officeId} [get]
func (h *Handler) GetByOfficeID(c *gin.Context) {
	officeID, err := strconv.ParseUint(c.Param("officeId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid office id")
		return
	}

	notaries, err := h.service.GetByOfficeID(uint(officeID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, notaries)
}

// GetByID godoc
// @Summary      Get notary by ID
// @Tags         Notaries
// @Produce      json
// @Param        id   path      int  true  "Notary ID"
// @Success      200  {object}  utils.APIResponse{data=Notary}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/notaries/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	n, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, n)
}

// Update godoc
// @Summary      Update a notary
// @Tags         Notaries
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Notary ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=Notary}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/notaries/{id} [put]
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

	n, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, n)
}

// Delete godoc
// @Summary      Delete a notary
// @Tags         Notaries
// @Produce      json
// @Param        id   path      int  true  "Notary ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/notaries/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "notary deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	notaries := r.Group("/notaries")
	{
		notaries.POST("", h.Create)
		notaries.GET("", h.GetAll)
		notaries.GET("/:id", h.GetByID)
		notaries.GET("/office/:officeId", h.GetByOfficeID)
		notaries.PUT("/:id", h.Update)
		notaries.DELETE("/:id", h.Delete)
	}
}
