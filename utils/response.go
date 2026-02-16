package utils

import (
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationRequest struct {
	Page    int `form:"page,default=1" binding:"min=1"`
	PerPage int `form:"per_page,default=50" binding:"min=1,max=500"`
}

type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

func Paginate(p PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (p.Page - 1) * p.PerPage
		return db.Offset(offset).Limit(p.PerPage)
	}
}

func NewPaginatedResponse(items interface{}, total int64, p PaginationRequest) PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(p.PerPage)))
	return PaginatedResponse{
		Items:      items,
		Total:      total,
		Page:       p.Page,
		PerPage:    p.PerPage,
		TotalPages: totalPages,
	}
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
	})
}
