package permission

import (
	"errors"
	"strings"

	"nfs-api/database"
	"nfs-api/utils"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(req CreateRequest) (*Permission, error) {
	perm := Permission{
		Name: strings.ToUpper(req.Name),
	}

	if err := database.DB.Create(&perm).Error; err != nil {
		return nil, errors.New("failed to create permission")
	}
	return &perm, nil
}

func (s *Service) GetAll(p utils.PaginationRequest) (*utils.PaginatedResponse, error) {
	var perms []Permission
	var total int64

	query := database.DB.Model(&Permission{}).Where("deleted = ?", 0)
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count permissions")
	}

	if err := query.Scopes(utils.Paginate(p)).Find(&perms).Error; err != nil {
		return nil, errors.New("failed to fetch permissions")
	}

	result := utils.NewPaginatedResponse(perms, total, p)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*Permission, error) {
	var perm Permission
	if err := database.DB.Where("deleted = ?", 0).First(&perm, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("permission not found")
		}
		return nil, errors.New("failed to fetch permission")
	}
	return &perm, nil
}

func (s *Service) Update(id uint, req UpdateRequest) (*Permission, error) {
	perm, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		perm.Name = *req.Name
	}

	if err := database.DB.Save(perm).Error; err != nil {
		return nil, errors.New("failed to update permission")
	}
	return perm, nil
}

func (s *Service) Delete(id uint) error {
	perm, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := database.DB.Model(perm).Update("deleted", 1).Error; err != nil {
		return errors.New("failed to delete permission")
	}
	return nil
}
