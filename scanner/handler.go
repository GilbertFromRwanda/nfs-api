package scanner

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

// Enroll godoc
// @Summary      Enroll a new scanner client
// @Tags         Scanner
// @Accept       json
// @Produce      json
// @Param        body  body      EnrollRequest  true  "Enroll payload"
// @Success      201   {object}  utils.APIResponse{data=ScannerClient}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/scanner/enroll [post]
// @Security BearerAuth
func (h *Handler) Enroll(c *gin.Context) {
	var req EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Inject office_id and user_id from JWT — never trust the request body for these
	officeID, _ := c.Get("office_id")
	if officeID == nil {
		utils.Error(c, http.StatusBadRequest, "office not found in token")
		return
	}
	userID, _ := c.Get("user_id")

	req.CompanyID = officeID.(uint)
	if userID != nil {
		req.UserID = userID.(uint)
	}

	client, err := h.service.Enroll(req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, client)
}

// List godoc
// @Summary      List scanner clients
// @Tags         Scanner
// @Produce      json
// @Param        q        query     string  false  "Search by name or NID"
// @Param        service  query     string  false  "Filter by service name"
// @Param        from     query     string  false  "From date (YYYY-MM-DD)"
// @Param        to       query     string  false  "To date (YYYY-MM-DD)"
// @Param        page     query     int     false  "Page number"  default(1)
// @Param        per_page query     int     false  "Items per page"  default(50)
// @Success      200  {object}  utils.APIResponse{data=utils.PaginatedResponse}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/scanner/clients [get]
// @Security BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Scope to the caller's office automatically
	if officeID, exists := c.Get("office_id"); exists {
		req.CompanyID = officeID.(uint)
	}

	result, err := h.service.List(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, result)
}

// UpdateFingerprint godoc
// @Summary      Update fingerprint template for an enrolled client
// @Tags         Scanner
// @Accept       json
// @Produce      json
// @Param        id    path      int                       true  "Client ID"
// @Param        body  body      UpdateFingerprintRequest  true  "New fingerprint"
// @Success      200   {object}  utils.APIResponse
// @Failure      400   {object}  utils.APIResponse
// @Failure      404   {object}  utils.APIResponse
// @Router       /api/v1/scanner/clients/{id}/fingerprint [patch]
// @Security BearerAuth
func (h *Handler) UpdateFingerprint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateFingerprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var officeID uint
	if oid, exists := c.Get("office_id"); exists {
		officeID = oid.(uint)
	}

	if err := h.service.UpdateFingerprint(uint(id), officeID, req); err != nil {
		if err.Error() == "not found" {
			utils.Error(c, http.StatusNotFound, "client not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to update fingerprint")
		}
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "fingerprint updated"})
}

// GetFingerprint godoc
// @Summary      Get fingerprint template for a client
// @Tags         Scanner
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  utils.APIResponse{data=string}
// @Failure      400  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/scanner/clients/{id}/fingerprint [get]
// @Security BearerAuth
func (h *Handler) GetFingerprint(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var officeID uint
	if oid, exists := c.Get("office_id"); exists {
		officeID = oid.(uint)
	}

	template, err := h.service.GetFingerprint(uint(id), officeID)
	if err != nil {
		if err.Error() == "not found" {
			utils.Error(c, http.StatusNotFound, "client not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to fetch fingerprint")
		}
		return
	}

	utils.Success(c, http.StatusOK, template)
}

// GetByID godoc
// @Summary      Get a scanner client by ID
// @Tags         Scanner
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  utils.APIResponse{data=ScannerClient}
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/scanner/clients/{id} [get]
// @Security BearerAuth
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var officeID uint
	if oid, exists := c.Get("office_id"); exists {
		officeID = oid.(uint)
	}

	client, err := h.service.Get(uint(id), officeID)
	if err != nil {
		if err.Error() == "not found" {
			utils.Error(c, http.StatusNotFound, "client not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to fetch client")
		}
		return
	}

	utils.Success(c, http.StatusOK, client)
}

// Delete godoc
// @Summary      Delete a scanner client
// @Tags         Scanner
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  utils.APIResponse
// @Failure      404  {object}  utils.APIResponse
// @Router       /api/v1/scanner/clients/{id} [delete]
// @Security BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var officeID uint
	if oid, exists := c.Get("office_id"); exists {
		officeID = oid.(uint)
	}

	if err := h.service.Delete(uint(id), officeID); err != nil {
		if err.Error() == "not found" {
			utils.Error(c, http.StatusNotFound, "client not found")
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to delete client")
		}
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "client deleted"})
}

// AllTemplates godoc
// @Summary      Get all fingerprint templates for the office (duplicate detection)
// @Tags         Scanner
// @Produce      json
// @Success      200  {object}  utils.APIResponse{data=[]TemplateEntry}
// @Failure      500  {object}  utils.APIResponse
// @Router       /api/v1/scanner/templates [get]
// @Security BearerAuth
func (h *Handler) AllTemplates(c *gin.Context) {
	var officeID uint
	if oid, exists := c.Get("office_id"); exists {
		officeID = oid.(uint)
	}

	entries, err := h.service.AllTemplates(officeID)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, entries)
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	sc := r.Group("/scanner")
	{
		sc.POST("/enroll", h.Enroll)
		sc.GET("/clients", h.List)
		sc.GET("/templates", h.AllTemplates)
		sc.GET("/clients/:id", h.GetByID)
		sc.GET("/clients/:id/fingerprint", h.GetFingerprint)
		sc.PATCH("/clients/:id/fingerprint", h.UpdateFingerprint)
		sc.DELETE("/clients/:id", h.Delete)
	}
}
