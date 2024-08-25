package repository

import (
	"context"
	"time"

	"github.com/billing-engine/internal/repository/entity"
	"github.com/billing-engine/internal/repository/models"
	"gorm.io/gorm"
)

type ILoanRepository interface {
	CreateLoan(ctx context.Context, data entity.LoanEntity) (entity.LoanEntity, error)
	Get(ctx context.Context, username string, status int) (entity.LoanEntity, error)
	UpdateStatus(ctx context.Context, loanId int, status int) error
	GetByStatus(ctx context.Context, status int) ([]entity.LoanEntity, error)
}

type LoanRepository struct {
	DB *gorm.DB
}

func NewLoanRepository(DB *gorm.DB) ILoanRepository {
	return &LoanRepository{
		DB: DB,
	}
}

func (lr *LoanRepository) UpdateStatus(ctx context.Context, loanId int, status int) error {
	model := models.LoanModel{
		Id: loanId,
	}

	if response := lr.DB.Table("loan").Model(&model).Updates(map[string]interface{}{
		"status": status,
	}); response.Error != nil {
		return response.Error
	}

	return nil
}

func (lr *LoanRepository) GetByStatus(ctx context.Context, status int) ([]entity.LoanEntity, error) {
	models := []models.LoanModel{}

	if response := lr.DB.Table("loan").Where("status = ?", status).Find(&models); response.Error != nil {
		return []entity.LoanEntity{}, response.Error
	}

	return convertBulkModelToEntitiesLoan(models), nil
}

func (lr *LoanRepository) CreateLoan(ctx context.Context, data entity.LoanEntity) (entity.LoanEntity, error) {
	model := models.LoanModel{
		Username:  data.Username,
		Amount:    data.Amount,
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"),
		Status:    data.Status,
	}
	if err := lr.DB.Table("loan").Create(&model); err.Error != nil {
		return entity.LoanEntity{}, err.Error
	}

	return convertModelToEntityLoan(model), nil
}

func (lr *LoanRepository) Get(ctx context.Context, username string, status int) (entity.LoanEntity, error) {
	model := models.LoanModel{}
	if response := lr.DB.Table("loan").Where("username = ?", username).Where("status  = ?", status).Last(&model); response.Error != nil {
		return entity.LoanEntity{}, response.Error
	}

	return convertModelToEntityLoan(model), nil
}

func convertModelToEntityLoan(model models.LoanModel) entity.LoanEntity {
	createdAt, _ := time.Parse("2006-01-02 15:04:05", model.CreatedAt)

	return entity.LoanEntity{
		Id:        model.Id,
		Username:  model.Username,
		Amount:    model.Amount,
		CreatedAt: createdAt,
		Status:    model.Status,
	}
}

func convertBulkModelToEntitiesLoan(models []models.LoanModel) []entity.LoanEntity {
	result := []entity.LoanEntity{}

	for _, model := range models {
		result = append(result, convertModelToEntityLoan(model))
	}

	return result
}
