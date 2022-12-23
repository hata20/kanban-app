package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var result entity.User
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", id).Find(&result)
	if err.Error != nil {
		return entity.User{}, err.Error
	}
	return result, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var result entity.User
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Find(&result)
	if err.Error != nil {
		return entity.User{}, err.Error
	}
	return result, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id= ?", user.ID).Updates(user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	err := r.db.WithContext(ctx).Delete(&entity.User{}, id)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
