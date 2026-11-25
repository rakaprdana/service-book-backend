package repositories

import (
	"context"
	"errors"
	commErr "service-book/common/errors"
	constErr "service-book/constants/errors"
	constUserErr "service-book/constants/errors/userErr"
	"service-book/domains/dto"
	models "service-book/domains/models/user"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(context.Context, *dto.RegisterRequest) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	FindByEmail(context.Context, string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Register(ctx context.Context, req *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:     uuid.New(),
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	}

	err := r.db.
		WithContext(ctx).
		Create(&user).
		Error

	if err != nil {
		return nil, commErr.WrapError(constErr.ErrSQLError)
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.
		WithContext(ctx).
		Preload("Role").
		Where("username = ? ", username).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constUserErr.ErrUserNotFound
		}

		return nil, commErr.WrapError(constErr.ErrSQLError)
	}

	return &user, nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := u.db.
		WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constUserErr.ErrUserNotFound
		}
		return nil, commErr.WrapError(constErr.ErrSQLError)
	}

	return &user, nil
}
