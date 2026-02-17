package officefeature

import (
	"errors"

	"nfs-api/database"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Assign(req AssignRequest) error {
	if len(req.AppFeatureIDs) == 0 {
		return errors.New("app feature IDs list cannot be empty")
	}

	var existing []OfficeFeature
	database.DB.Where("office_id = ? AND app_feature_id IN ?", req.OfficeID, req.AppFeatureIDs).Find(&existing)

	existingMap := make(map[uint]bool)
	for _, of := range existing {
		existingMap[of.AppFeatureID] = true
	}

	var newFeatures []OfficeFeature
	for _, featureID := range req.AppFeatureIDs {
		if !existingMap[featureID] {
			newFeatures = append(newFeatures, OfficeFeature{
				OfficeID:     req.OfficeID,
				AppFeatureID: featureID,
			})
		}
	}

	if len(newFeatures) == 0 {
		return errors.New("all features are already assigned to this office")
	}

	if err := database.DB.Create(&newFeatures).Error; err != nil {
		return errors.New("failed to assign features")
	}
	return nil
}

func (s *Service) Revoke(req RevokeRequest) error {
	if len(req.AppFeatureIDs) == 0 {
		return errors.New("app feature IDs list cannot be empty")
	}

	result := database.DB.Where("office_id = ? AND app_feature_id IN ?", req.OfficeID, req.AppFeatureIDs).Delete(&OfficeFeature{})
	if result.Error != nil {
		return errors.New("failed to revoke features")
	}
	if result.RowsAffected == 0 {
		return errors.New("no matching office features found")
	}
	return nil
}

func (s *Service) GetByOfficeID(officeID uint) ([]OfficeFeatureResponse, error) {
	var results []OfficeFeatureResponse
	err := database.DB.Table("office_features").
		Select("office_features.id, office_features.office_id, office_features.app_feature_id, app_features.name as app_feature_name, app_features.status").
		Joins("LEFT JOIN app_features ON app_features.id = office_features.app_feature_id").
		Where("office_features.office_id = ?", officeID).
		Scan(&results).Error
	if err != nil {
		return nil, errors.New("failed to fetch office features")
	}
	return results, nil
}
