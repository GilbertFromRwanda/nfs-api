package nlachecklist

import (
	"time"

	"nfs-api/auth"
	"nfs-api/office"

	"gorm.io/gorm"
)

type NlaChecklist struct {
	ID           uint                `gorm:"primaryKey" json:"id"`
	UPI          string              `gorm:"size:40;not null;index" json:"upi"`
	Seller1Names string              `gorm:"size:200;not null" json:"seller_1_names"`
	Seller2Names string              `gorm:"size:200" json:"seller_2_names"`
	Buyer1Names  string              `gorm:"size:200;not null" json:"buyer_1_names"`
	Buyer2Names  string              `gorm:"size:200" json:"buyer_2_names"`
	UserID       uint                `gorm:"not null;index" json:"user_id"`
	User         *auth.User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OfficeID     uint                `gorm:"not null;index" json:"office_id"`
	Office       *office.Office      `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
	Date         time.Time           `gorm:"not null;index" json:"date"`
	Dossiers     []ChecklistDossier  `gorm:"foreignKey:NlaChecklistID;constraint:OnDelete:CASCADE" json:"dossiers"`
	Status       string              `gorm:"-" json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    gorm.DeletedAt     `gorm:"index" json:"deleted_at,omitempty"`
}

type ChecklistDossier struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	NlaChecklistID uint   `gorm:"not null;index" json:"nla_checklist_id"`
	Name           string `gorm:"size:200;not null" json:"name"`
	Status         string `gorm:"size:20;not null;default:'no'" json:"status"` // yes | no | no_needed
}

func (c *NlaChecklist) ComputeStatus() {
	c.Status = "completed"
	for _, d := range c.Dossiers {
		if d.Status == "no" {
			c.Status = "uncompleted"
			return
		}
	}
}
