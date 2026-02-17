package officefeature

type AssignRequest struct {
	OfficeID      uint   `json:"office_id" binding:"required"`
	AppFeatureIDs []uint `json:"app_feature_ids" binding:"required,min=1"`
}

type RevokeRequest struct {
	OfficeID      uint   `json:"office_id" binding:"required"`
	AppFeatureIDs []uint `json:"app_feature_ids" binding:"required,min=1"`
}

type OfficeFeatureResponse struct {
	ID             uint   `json:"id"`
	OfficeID       uint   `json:"office_id"`
	AppFeatureID   uint   `json:"app_feature_id"`
	AppFeatureName string `json:"app_feature_name"`
	Status         string `json:"status"`
}
