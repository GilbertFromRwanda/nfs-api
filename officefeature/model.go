package officefeature

import "time"

type OfficeFeature struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OfficeID     uint      `gorm:"not null;index;uniqueIndex:idx_office_feature" json:"office_id"`
	AppFeatureID uint      `gorm:"not null;index;uniqueIndex:idx_office_feature" json:"app_feature_id"`
	CreatedAt    time.Time `json:"created_at"`
}
