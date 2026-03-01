package agreement

import (
	"net/http"
	"strconv"

	"nfs-api/utils"

	"github.com/gin-gonic/gin"
)

// ListClients godoc
// @Summary      List clients linked to an agreement
// @Tags         Agreements
// @Produce      json
// @Param        id   path      int  true  "Agreement ID"
// @Success      200  {object}  utils.APIResponse{data=[]AgreementClient}
// @Failure      400  {object}  utils.APIResponse
// @Router       /api/v1/agreements/{id}/clients [get]
// @Security BearerAuth
func (h *Handler) ListClients(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	clients, err := h.service.ListClients(uint(id))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, clients)
}

// AddClient godoc
// @Summary      Link a scanner client to an agreement
// @Tags         Agreements
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "Agreement ID"
// @Param        body  body      AddClientRequest  true  "Client payload"
// @Success      201   {object}  utils.APIResponse{data=AgreementClient}
// @Failure      400   {object}  utils.APIResponse
// @Router       /api/v1/agreements/{id}/clients [post]
// @Security BearerAuth
func (h *Handler) AddClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req AddClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var officeID uint
	if oid, exists := c.Get("office_id"); exists {
		officeID = oid.(uint)
	}

	ac, err := h.service.AddClient(uint(id), officeID, req)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, ac)
}

// RemoveClient godoc
// @Summary      Remove a client link from an agreement
// @Tags         Agreements
// @Produce      json
// @Param        id      path      int  true  "Agreement ID"
// @Param        linkID  path      int  true  "AgreementClient link ID"
// @Success      200     {object}  utils.APIResponse
// @Failure      400     {object}  utils.APIResponse
// @Failure      404     {object}  utils.APIResponse
// @Router       /api/v1/agreements/{id}/clients/{linkID} [delete]
// @Security BearerAuth
func (h *Handler) RemoveClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	linkID, err := strconv.ParseUint(c.Param("linkID"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid link id")
		return
	}

	var officeID uint
	if oid, exists := c.Get("office_id"); exists {
		officeID = oid.(uint)
	}

	if err := h.service.RemoveClient(uint(id), officeID, uint(linkID)); err != nil {
		utils.Error(c, http.StatusNotFound, err.Error())
		return
	}

	utils.Success(c, http.StatusOK, gin.H{"message": "client removed"})
}
