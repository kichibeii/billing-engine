package repository

import (
	"context"
	"time"

	"github.com/billing-engine/internal/commons"
	"github.com/billing-engine/internal/repository/entity"
	"github.com/billing-engine/internal/repository/models"
	"gorm.io/gorm"
)

type IPayLoanRepository interface {
	GetPayLoanByLoanId(ctx context.Context, loandId int) ([]entity.PayLoanEntity, error)
	GetInSpecificTimeAndStatus(ctx context.Context, loanId int, timeNow time.Time) ([]entity.PayLoanEntity, error)
	BatchInsert(ctx context.Context, datas []entity.PayLoanEntity) error
	Update(ctx context.Context, id int, data entity.PayLoanEntity) error
}

type PayLoanRepository struct {
	DB *gorm.DB
}

func NewPayLoanRepository(DB *gorm.DB) IPayLoanRepository {
	return &PayLoanRepository{
		DB: DB,
	}
}

func (plr *PayLoanRepository) GetInSpecificTimeAndStatus(ctx context.Context, loanId int, timeNow time.Time) ([]entity.PayLoanEntity, error) {
	models := []models.PayLoanModel{}

	if response := plr.DB.Table("pay_loan").
		Where("loan_id = ?", loanId).
		Where("status = ?", commons.StatusPayLoanUnpayed).
		Where("created_at < ?", timeNow.Format("2006-01-02 15:04:05")).
		Find(&models); response.Error != nil {
		return []entity.PayLoanEntity{}, response.Error
	}

	return convertBulkModelToEntitiesPayLoan(models), nil
}

func (plr *PayLoanRepository) BatchInsert(ctx context.Context, datas []entity.PayLoanEntity) error {
	models := convertBulkEntityToModelsPayLoan(datas)

	if response := plr.DB.Table("pay_loan").Create(&models); response.Error != nil {
		return response.Error
	}

	return nil
}

func (plr *PayLoanRepository) Update(ctx context.Context, id int, data entity.PayLoanEntity) error {
	model := models.PayLoanModel{
		Id: id,
	}
	if response := plr.DB.Table("pay_loan").Model(&model).Updates(map[string]interface{}{
		"status": data.Status,
	}); response.Error != nil {
		return response.Error
	}

	return nil
}

func (plr *PayLoanRepository) GetPayLoanByLoanId(ctx context.Context, loandId int) ([]entity.PayLoanEntity, error) {
	models := []models.PayLoanModel{}

	if response := plr.DB.Table("pay_loan").Where("loan_id = ?", loandId).Scan(&models); response.Error != nil {
		return []entity.PayLoanEntity{}, response.Error
	}

	return convertBulkModelToEntitiesPayLoan(models), nil
}

func convertModelToEntityPayLoan(model models.PayLoanModel) entity.PayLoanEntity {
	createdAt, _ := time.Parse("2006-01-02 15:04:05", model.CreatedAt)

	return entity.PayLoanEntity{
		Id:        model.Id,
		LoanId:    model.LoanId,
		Amount:    model.Amount,
		Status:    model.Status,
		CreatedAt: createdAt,
	}
}

func convertEntityToModelPayLoan(entity entity.PayLoanEntity) models.PayLoanModel {
	return models.PayLoanModel{
		Id:        entity.Id,
		LoanId:    entity.LoanId,
		Amount:    entity.Amount,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
		Status:    entity.Status,
	}
}

func convertBulkEntityToModelsPayLoan(entities []entity.PayLoanEntity) []models.PayLoanModel {
	result := []models.PayLoanModel{}

	for _, entity := range entities {
		result = append(result, convertEntityToModelPayLoan(entity))
	}

	return result
}

func convertBulkModelToEntitiesPayLoan(models []models.PayLoanModel) []entity.PayLoanEntity {
	result := []entity.PayLoanEntity{}

	for _, model := range models {
		result = append(result, convertModelToEntityPayLoan(model))
	}

	return result
}
