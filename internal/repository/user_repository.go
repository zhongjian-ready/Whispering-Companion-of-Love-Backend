package repository

import (
	"miniapp-backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindAllBasicInfo() ([]model.User, int64, error) {
	var users []model.User
	var count int64

	// Get total count
	result := r.db.Model(&model.User{}).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Get basic info list
	// Only selecting username, nickname, gender as requested
	result = r.db.Select("username", "nickname", "gender").Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, count, nil
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) UpdateSettings(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepository) FindByOpenID(openid string) (*model.User, error) {
	var user model.User
	err := r.db.Where("openid = ?", openid).First(&user).Error
	return &user, err
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
