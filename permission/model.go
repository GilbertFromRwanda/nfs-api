package permission

import "time"

type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null,index" json:"name"`
	Deleted   int       `gorm:"default:0" json:"deleted"` // 0 = active, 1 = deleted
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
