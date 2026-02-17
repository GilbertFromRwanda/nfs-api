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

// Assign godoc
// @Summary      Assign a permission to a user
// @Tags         User Permissions
// @Accept       json
// @Produce      json
// @Param        body  body      AssignRequest  true  "Assign payload"
// @Success      201   {object}  utils.APIResponse{data=UserPermissionResponse}
// @Failure      400   {object}  utils.APIResponse
// @Failure      409   {object}  utils.APIResponse
// @Router       /api/user-permissions [post]
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

// GetByUserID godoc
// @Summary      Get permissions by user ID
// @Tags         User Permissions
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Success      200     {object}  utils.APIResponse{data=[]UserPermissionResponse}
// @Failure      400     {object}  utils.APIResponse
// @Router       /api/user-permissions/user/{userId} [get]
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

// Revoke godoc
// @Summary      Revoke a user permission
// @Tags         User Permissions
// @Produce      json
// @Param        id   path      int  true  "User Permission ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/user-permissions/{id} [delete]
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

// RevokeMany godoc
// @Summary      Revoke multiple permissions from a user
// @Tags         User Permissions
// @Accept       json
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Param        body    body      object{permission_ids=[]int}  true  "Permission IDs"
// @Success      200     {object}  utils.APIResponse
// @Failure      400     {object}  utils.APIResponse
// @Failure      404     {object}  utils.APIResponse
// @Router       /api/user-permissions/user/{userId}/revoke-many [delete]
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

// AssignMany godoc
// @Summary      Assign multiple permissions to a user
// @Tags         User Permissions
// @Accept       json
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Param        body    body      object{permission_ids=[]int}  true  "Permission IDs"
// @Success      201     {object}  utils.APIResponse
// @Failure      400     {object}  utils.APIResponse
// @Failure      409     {object}  utils.APIResponse
// @Router       /api/user-permissions/user/{userId}/assign-many [post]
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
