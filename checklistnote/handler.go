package checklistnote

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
// @Summary      Create a new checklist note
// @Tags         Checklist Notes
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=ChecklistNote}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/checklist-notes [post]
// @Security BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	note, err := h.service.Create(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, note)
}

// GetAll godoc
// @Summary      List all checklist notes
// @Tags         Checklist Notes
// @Produce      json
// @Param        page       query     int  false  "Page number"  default(1)
// @Param        per_page   query     int  false  "Items per page"  default(50)
// @Param        office_id  query     int  false  "Filter by office ID"
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/checklist-notes [get]
// @Security BearerAuth
func (h *Handler) GetAll(c *gin.Context) {
	var filter FilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	notes, err := h.service.GetAll(filter)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, notes)
}

// GetByID godoc
// @Summary      Get checklist note by ID
// @Tags         Checklist Notes
// @Produce      json
// @Param        id   path      int  true  "Checklist Note ID"
// @Success      200  {object}  utils.APIResponse{data=ChecklistNote}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/checklist-notes/{id} [get]
// @Security BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	note, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, note)
}

// Update godoc
// @Summary      Update a checklist note
// @Tags         Checklist Notes
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "Checklist Note ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=ChecklistNote}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/checklist-notes/{id} [put]
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

	note, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, note)
}

// Delete godoc
// @Summary      Delete a checklist note
// @Tags         Checklist Notes
// @Produce      json
// @Param        id   path      int  true  "Checklist Note ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/checklist-notes/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "checklist note deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	notes := r.Group("/checklist-notes")
	{
		notes.POST("", h.Create)
		notes.GET("", h.GetAll)
		notes.GET("/:id", h.GetByID)
		notes.PUT("/:id", h.Update)
		notes.DELETE("/:id", h.Delete)
	}
}
