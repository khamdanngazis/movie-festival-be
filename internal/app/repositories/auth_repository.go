package repositories

import (
	"movie-festival-be/internal/app/entities"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByEmail(username string) (*entities.User, error)
	FindUserByID(id uint) (*entities.User, error)
	SaveUser(user *entities.User) error
	UpdateLoginStatus(userID uint, status bool) error
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) FindUserByID(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepositoryImpl) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepositoryImpl) SaveUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepositoryImpl) UpdateLoginStatus(userID uint, status bool) error {
	return r.db.Model(&entities.User{}).Where("id = ?", userID).Update("logged_in", status).Error
}
