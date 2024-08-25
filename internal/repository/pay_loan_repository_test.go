package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/billing-engine/internal/commons"
	"github.com/billing-engine/internal/repository/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to open gorm DB: %v", err)
	}

	return db, mock
}

func TestPayLoanRepository_GetInSpecificTimeAndStatus(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewPayLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		loanId := 1
		timeNow := time.Now()

		rows := sqlmock.NewRows([]string{"id", "loan_id", "amount", "status", "created_at"}).
			AddRow(1, loanId, 1000.0, commons.StatusPayLoanUnpayed, "2023-08-24 10:00:00")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `pay_loan` WHERE loan_id = ? AND status = ? AND created_at < ?")).
			WithArgs(loanId, commons.StatusPayLoanUnpayed, timeNow.Format("2006-01-02 15:04:05")).
			WillReturnRows(rows)

		results, err := repo.GetInSpecificTimeAndStatus(context.Background(), loanId, timeNow)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, loanId, results[0].LoanId)
		assert.Equal(t, commons.StatusPayLoanUnpayed, results[0].Status)
	})

	t.Run("error", func(t *testing.T) {
		loanId := 1
		timeNow := time.Now()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `pay_loan` WHERE loan_id = ? AND status = ? AND created_at < ?")).
			WithArgs(loanId, commons.StatusPayLoanUnpayed, timeNow.Format("2006-01-02 15:04:05")).
			WillReturnError(gorm.ErrRecordNotFound)

		results, err := repo.GetInSpecificTimeAndStatus(context.Background(), loanId, timeNow)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}

func TestPayLoanRepository_BatchInsert(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewPayLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		data := []entity.PayLoanEntity{
			{LoanId: 1, Amount: 1000.0, Status: commons.StatusPayLoanUnpayed, CreatedAt: time.Now()},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `pay_loan` (`loan_id`,`amount`,`created_at`,`status`) VALUES (?,?,?,?)")).
			WithArgs(1, 1000.0, sqlmock.AnyArg(), commons.StatusPayLoanUnpayed).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.BatchInsert(context.Background(), data)

		assert.NoError(t, err)
		mock.ExpectationsWereMet()
	})

	t.Run("error", func(t *testing.T) {
		data := []entity.PayLoanEntity{
			{LoanId: 1, Amount: 1000.0, Status: commons.StatusPayLoanUnpayed, CreatedAt: time.Now()},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `pay_loan` (`loan_id`,`amount`,`created_at`,`status`) VALUES (?,?,?,?)")).
			WithArgs(1, 1000.0, sqlmock.AnyArg(), commons.StatusPayLoanUnpayed).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.BatchInsert(context.Background(), data)

		assert.Error(t, err)
		mock.ExpectationsWereMet()
	})
}

func TestPayLoanRepository_Update(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewPayLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		data := entity.PayLoanEntity{Status: commons.StatusPayLoanUnpayed}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `pay_loan` SET `status`=? WHERE `id` = ?")).
			WithArgs(data.Status, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), 1, data)

		assert.NoError(t, err)
		mock.ExpectationsWereMet()
	})

	t.Run("error", func(t *testing.T) {
		data := entity.PayLoanEntity{Status: commons.StatusPayLoanUnpayed}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `pay_loan` SET `status`=? WHERE `id` = ?")).
			WithArgs(data.Status, 1).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.Update(context.Background(), 1, data)

		assert.Error(t, err)
		mock.ExpectationsWereMet()
	})
}

func TestPayLoanRepository_GetPayLoanByLoanId(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewPayLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		loanId := 1

		rows := sqlmock.NewRows([]string{"id", "loan_id", "amount", "status", "created_at"}).
			AddRow(1, loanId, 1000.0, commons.StatusPayLoanUnpayed, "2023-08-24 10:00:00")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `pay_loan` WHERE loan_id = ?")).
			WithArgs(loanId).
			WillReturnRows(rows)

		results, err := repo.GetPayLoanByLoanId(context.Background(), loanId)

		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, loanId, results[0].LoanId)
		assert.Equal(t, commons.StatusPayLoanUnpayed, results[0].Status)
	})

	t.Run("error", func(t *testing.T) {
		loanId := 1

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `pay_loan` WHERE loan_id = ?")).
			WithArgs(loanId).
			WillReturnError(gorm.ErrRecordNotFound)

		results, err := repo.GetPayLoanByLoanId(context.Background(), loanId)

		assert.Error(t, err)
		assert.Empty(t, results)
	})
}
