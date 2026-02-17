package checklistnote

import (
	"time"

	"nfs-api/office"
)

type ChecklistNote struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:longtext;not null" json:"name"`
	OfficeID  uint           `gorm:"not null;index" json:"office_id"`
	Office    *office.Office `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
