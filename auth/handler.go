package auth

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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Register(req)
	if err != nil {
		utils.Error(c, http.StatusConflict, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, res)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.service.Login(req)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, res)
}

func (h *Handler) GetAll(c *gin.Context) {
	var p utils.PaginationRequest
	if err := c.ShouldBindQuery(&p); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.service.GetAll(p)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, users)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, user)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, user)
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

	utils.Success(c, http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)

	users := r.Group("/users")
	{
		users.GET("", h.GetAll)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}
