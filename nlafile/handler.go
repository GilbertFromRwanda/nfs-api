package nlafile

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
// @Summary      Create a new NLA file
// @Tags         NLA Files
// @Accept       json
// @Produce      json
// @Param        body  body      CreateRequest  true  "Create payload"
// @Success      201   {object}  utils.APIResponse{data=NlaFile}
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Router       /api/nla-files [post]
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

	nla, err := h.service.Create(req, userID.(uint))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, nla)
}

// GetAll godoc
// @Summary      List all NLA files
// @Tags         NLA Files
// @Produce      json
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        per_page   query     int     false  "Items per page"  default(50)
// @Param        upi        query     string  false  "Filter by UPI"
// @Param        name       query     string  false  "Filter by seller/buyer name"
// @Param        office_id  query     int     false  "Filter by office ID"
// @Param        status     query     string  false  "Filter by status"
// @Param        date_from  query     string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to    query     string  false  "Filter to date (YYYY-MM-DD)"
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/nla-files [get]
func (h *Handler) GetAll(c *gin.Context) {
	var filter FilterRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	files, err := h.service.GetAll(filter)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, files)
}

// GetByID godoc
// @Summary      Get NLA file by ID
// @Tags         NLA Files
// @Produce      json
// @Param        id   path      int  true  "NLA File ID"
// @Success      200  {object}  utils.APIResponse{data=NlaFile}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/nla-files/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	nla, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, nla)
}

// Update godoc
// @Summary      Update an NLA file
// @Tags         NLA Files
// @Accept       json
// @Produce      json
// @Param        id    path      int            true  "NLA File ID"
// @Param        body  body      UpdateRequest  true  "Update payload"
// @Success      200   {object}  utils.APIResponse{data=NlaFile}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/nla-files/{id} [put]
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

	nla, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, nla)
}

// Delete godoc
// @Summary      Delete an NLA file
// @Tags         NLA Files
// @Produce      json
// @Param        id   path      int  true  "NLA File ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/nla-files/{id} [delete]
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

	utils.Success(c, http.StatusOK, gin.H{"message": "nla file deleted"})
}

// MoveToNext godoc
// @Summary      Move NLA file to next status
// @Tags         NLA Files
// @Produce      json
// @Param        id   path      int  true  "NLA File ID"
// @Success      200  {object}  utils.APIResponse{data=NlaFile}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/nla-files/{id}/next [patch]
func (h *Handler) MoveToNext(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	nla, err := h.service.MoveToNext(uint(id))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, nla)
}

// Approve godoc
// @Summary      Approve an NLA file
// @Tags         NLA Files
// @Produce      json
// @Param        id   path      int  true  "NLA File ID"
// @Success      200  {object}  utils.APIResponse{data=NlaFile}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/nla-files/{id}/approve [patch]
func (h *Handler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	nla, err := h.service.Approve(uint(id))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, nla)
}

// Reject godoc
// @Summary      Reject an NLA file
// @Tags         NLA Files
// @Produce      json
// @Param        id   path      int  true  "NLA File ID"
// @Success      200  {object}  utils.APIResponse{data=NlaFile}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/nla-files/{id}/reject [patch]
func (h *Handler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	nla, err := h.service.Reject(uint(id))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, nla)
}

// GetStatusCounts godoc
// @Summary      Get NLA file status counts
// @Tags         NLA Files
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]StatusCount}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/nla-files/counts [get]
func (h *Handler) GetStatusCounts(c *gin.Context) {
	counts, err := h.service.GetStatusCounts()
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, counts)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	nla := r.Group("/nla-files")
	{
		nla.POST("", h.Create)
		nla.GET("", h.GetAll)
		nla.GET("/counts", h.GetStatusCounts)
		nla.GET("/:id", h.GetByID)
		nla.PUT("/:id", h.Update)
		nla.DELETE("/:id", h.Delete)
		nla.PATCH("/:id/next", h.MoveToNext)
		nla.PATCH("/:id/approve", h.Approve)
		nla.PATCH("/:id/reject", h.Reject)
	}
}
