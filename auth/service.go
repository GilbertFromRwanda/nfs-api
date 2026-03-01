package auth

import (
	"errors"
	"time"

	"nfs-api/config"
	"nfs-api/database"
	"nfs-api/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) Register(req RegisterRequest) (*AuthResponse, error) {
	var existing User
	if err := database.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  string(hashedPassword),
		OfficeID:  req.OfficeID,
		Role:      req.Role,
		Status:    status,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  toUserInfo(user),
	}, nil
}

// InviteStaff creates an employee account under the office_admin's office.
// officeID is taken from JWT — never from the request body.
func (s *Service) InviteStaff(req InviteStaffRequest, officeID uint) (*User, error) {
	var existing User
	if err := database.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return nil, errors.New("username already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	role := "employee"
	if req.Role != "" {
		role = req.Role
	}

	user := User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Password:  string(hashed),
		OfficeID:  &officeID,
		Role:      &role,
		Status:    "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create staff account")
	}

	return &user, nil
}

func (s *Service) Login(req LoginRequest) (*AuthResponse, error) {
	var user User
	if err := database.DB.Preload("Office").Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, errors.New("failed to find user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	if user.Status != "active" {
		return nil, errors.New("account is inactive")
	}
	if user.Office != nil && user.Office.Status != "active" {
		return nil, errors.New("associated office is inactive")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  toUserInfo(user),
	}, nil
}

func (s *Service) GetAll(p utils.PaginationRequest) (*utils.PaginatedResponse, error) {
	var users []User
	var total int64

	query := database.DB.Model(&User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count users")
	}

	if err := query.
		Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).
		Scopes(utils.Paginate(p)).
		Find(&users).Error; err != nil {
		return nil, errors.New("failed to fetch users")
	}

	result := utils.NewPaginatedResponse(users, total, p)
	return &result, nil
}

func (s *Service) GetByID(id uint) (*User, error) {
	var user User
	if err := database.DB.Preload("Office", func(db *gorm.DB) *gorm.DB { return db.Select("id, name, location, status") }).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to fetch user")
	}
	return &user, nil
}

func (s *Service) Update(id uint, req UpdateUserRequest) (*User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Username != nil {
		var existing User
		if err := database.DB.Where("username = ? AND id != ?", *req.Username, id).First(&existing).Error; err == nil {
			return nil, errors.New("username already taken")
		}
		updates["username"] = *req.Username
	}
	if req.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		updates["password"] = string(hashed)
	}
	if req.OfficeID != nil {
		updates["office_id"] = *req.OfficeID
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(user).Updates(updates).Error; err != nil {
			return nil, errors.New("failed to update user")
		}
	}

	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	result := database.DB.Delete(&User{}, id)
	if result.Error != nil {
		return errors.New("failed to delete user")
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (s *Service) ResetPassword(id uint, req ResetPasswordRequest) error {
	var user User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("failed to find user")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := database.DB.Model(&user).Update("password", string(hashed)).Error; err != nil {
		return errors.New("failed to reset password")
	}

	return nil
}

func (s *Service) generateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	if user.Role != nil {
		claims["role"] = *user.Role
	}

	if user.OfficeID != nil {
		claims["office_id"] = *user.OfficeID
		// Load office name if not already preloaded
		if user.Office != nil {
			claims["office_name"] = user.Office.Name
		} else {
			var office struct{ Name string }
			if err := database.DB.Table("offices").Select("name").Where("id = ?", *user.OfficeID).Scan(&office).Error; err == nil {
				claims["office_name"] = office.Name
			}
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func toUserInfo(user User) UserInfo {
	officeName := ""
	if user.Office != nil {
		officeName = user.Office.Name
	}
	return UserInfo{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Username:   user.Username,
		OfficeID:   user.OfficeID,
		OfficeName: officeName,
		Role:       user.Role,
		Status:     user.Status,
	}
}
