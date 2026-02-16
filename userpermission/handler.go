package userpermission

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

func (h *Handler) Assign(c *gin.Context) {
	var req AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	up, err := h.service.Assign(req)
	if err != nil {
		utils.Error(c, http.StatusConflict, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, up)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	perms, err := h.service.GetByUserID(uint(userID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, perms)
}

func (h *Handler) Revoke(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Revoke(uint(id)); err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "permission revoked"})
}

func (h *Handler) RevokeMany(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.RevokeMany(uint(userID), req.PermissionIDs); err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "permissions revoked"})
}
func (h *Handler) AssignMany(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AssignMany(uint(userID), req.PermissionIDs); err != nil {
		utils.Error(c, http.StatusConflict, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, gin.H{"message": "permissions assigned"})
}
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	up := r.Group("/user-permissions")
	{
		up.POST("", h.Assign)
		up.GET("/user/:userId", h.GetByUserID)
		up.DELETE("/:id", h.Revoke)
		up.POST("/user/:userId/assign-many", h.AssignMany)
		up.DELETE("/user/:userId/revoke-many", h.RevokeMany)
	}
}
