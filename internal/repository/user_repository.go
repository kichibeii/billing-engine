package repository

import (
	"context"
	"time"

	"github.com/billing-engine/internal/commons"
	"github.com/billing-engine/internal/repository/entity"
	"github.com/billing-engine/internal/repository/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUser(ctx context.Context, username string) (entity.UserEntity, error)
	Create(ctx context.Context, username string) (entity.UserEntity, error)
	UpdateUser(ctx context.Context, username string, status int) error
}

type UserRepository struct {
	DB *gorm.DB
}

type CreateUserModel struct {
	Username  string    `db:"username"`
	Amount    float64   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUserRepository(DB *gorm.DB) IUserRepository {
	return &UserRepository{
		DB: DB,
	}
}

func (ur *UserRepository) Create(ctx context.Context, username string) (entity.UserEntity, error) {
	model := models.UserModel{
		Username: username,
		Status:   commons.StatusUserNew,
	}

	if response := ur.DB.Table("users").Create(&model); response.Error != nil {
		return entity.UserEntity{}, response.Error
	}

	return entity.UserEntity{
		Username: model.Username,
		Status:   model.Status,
	}, nil
}

func (ur *UserRepository) GetUser(ctx context.Context, username string) (entity.UserEntity, error) {
	model := models.UserModel{}

	if response := ur.DB.Table("users").Where("username = ?", username).Find(&model); response.Error != nil {
		return entity.UserEntity{}, response.Error
	}

	return entity.UserEntity{
		Username: model.Username,
		Status:   model.Status,
	}, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, username string, status int) error {
	if response := ur.DB.Table("users").Where("username = ?", username).Update("status", status); response.Error != nil {
		return response.Error
	}

	return nil
}
