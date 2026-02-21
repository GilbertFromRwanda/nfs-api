package officefeature

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
// @Summary      Assign features to an office
// @Tags         Office Features
// @Accept       json
// @Produce      json
// @Param        body  body      AssignRequest  true  "Assign payload"
// @Success      201   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      409   {object}  utils.APIResponse
// @Router       /api/v1/office-features/assign [post]
func (h *Handler) Assign(c *gin.Context) {
	var req AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Assign(req); err != nil {
		utils.Error(c, http.StatusConflict, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, gin.H{"message": "features assigned"})
}

// Revoke godoc
// @Summary      Revoke features from an office
// @Tags         Office Features
// @Accept       json
// @Produce      json
// @Param        body  body      RevokeRequest  true  "Revoke payload"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      404   {object}  utils.APIResponse
// @Router       /api/v1/office-features/revoke [delete]
func (h *Handler) Revoke(c *gin.Context) {
	var req RevokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Revoke(req); err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "features revoked"})
}

// GetByOfficeID godoc
// @Summary      Get features assigned to an office
// @Tags         Office Features
// @Produce      json
// @Param        officeId  path      int  true  "Office ID"
// @Success      200       {object}  utils.APIResponse{data=[]OfficeFeatureResponse}
// @Failure      400       {object}  utils.APIResponse
// @Router       /api/v1/office-features/office/{officeId} [get]
func (h *Handler) GetByOfficeID(c *gin.Context) {
	officeID, err := strconv.ParseUint(c.Param("officeId"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid office id")
		return
	}

	features, err := h.service.GetByOfficeID(uint(officeID))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, features)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	of := r.Group("/office-features")
	{
		of.POST("/assign", h.Assign)
		of.DELETE("/revoke", h.Revoke)
		of.GET("/office/:officeId", h.GetByOfficeID)
	}
}
