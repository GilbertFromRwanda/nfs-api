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
