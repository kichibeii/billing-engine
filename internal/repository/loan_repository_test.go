package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/billing-engine/internal/repository/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLoanRepository_CreateLoan(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		data := entity.LoanEntity{
			Username:  "user123",
			Amount:    1000.0,
			CreatedAt: time.Now(),
			Status:    1,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `loan` (`username`,`amount`,`created_at`,`status`) VALUES (?,?,?,?)")).
			WithArgs(data.Username, data.Amount, data.CreatedAt.Format("2006-01-02 15:04:05"), data.Status).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result, err := repo.CreateLoan(context.Background(), data)

		assert.NoError(t, err)
		assert.Equal(t, data.Username, result.Username)
		assert.Equal(t, data.Status, result.Status)
	})

	t.Run("error", func(t *testing.T) {
		data := entity.LoanEntity{
			Username:  "user123",
			Amount:    1000.0,
			CreatedAt: time.Now(),
			Status:    1,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `loan` (`username`,`amount`,`created_at`,`status`) VALUES (?,?,?,?)")).
			WithArgs(data.Username, data.Amount, data.CreatedAt.Format("2006-01-02 15:04:05"), data.Status).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		_, err := repo.CreateLoan(context.Background(), data)

		assert.Error(t, err)
	})
}

func TestLoanRepository_Get(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		username := "user123"
		status := 1

		row := sqlmock.NewRows([]string{"id", "username", "amount", "status", "created_at"}).
			AddRow(1, username, 1000.0, status, "2023-08-24 10:00:00")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `loan` WHERE username = ? AND status = ? ORDER BY `loan`.`id` DESC LIMIT ?")).
			WithArgs(username, status, 1).
			WillReturnRows(row)

		result, err := repo.Get(context.Background(), username, status)

		assert.NoError(t, err)
		assert.Equal(t, username, result.Username)
		assert.Equal(t, status, result.Status)
	})

	t.Run("error", func(t *testing.T) {
		username := "user123"
		status := 1

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `loan` WHERE username = ? AND status = ? ORDER BY `loan`.`id` DESC LIMIT ?")).
			WithArgs(username, status, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.Get(context.Background(), username, status)

		assert.Error(t, err)
	})
}

func TestLoanRepository_UpdateStatus(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		loanId := 1
		status := 2

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `loan` SET `status`=? WHERE `id` = ?")).
			WithArgs(status, loanId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.UpdateStatus(context.Background(), loanId, status)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		loanId := 1
		status := 2

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `loan` SET `status`=? WHERE `id` = ?")).
			WithArgs(status, loanId).
			WillReturnError(gorm.ErrInvalidData)
		mock.ExpectRollback()

		err := repo.UpdateStatus(context.Background(), loanId, status)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestLoanRepository_GetByStatus(t *testing.T) {
	db, mock := setupTestDB(t)

	repo := NewLoanRepository(db)

	t.Run("success", func(t *testing.T) {
		status := 1

		rows := sqlmock.NewRows([]string{"id", "username", "amount", "status", "created_at"}).
			AddRow(1, "user123", 1000.0, status, "2023-08-24 10:00:00").
			AddRow(2, "user456", 2000.0, status, "2023-08-24 11:00:00")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `loan` WHERE status = ?")).
			WithArgs(status).
			WillReturnRows(rows)

		results, err := repo.GetByStatus(context.Background(), status)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "user123", results[0].Username)
		assert.Equal(t, "user456", results[1].Username)
	})

	t.Run("error", func(t *testing.T) {
		status := 1

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `loan` WHERE status = ?")).
			WithArgs(status).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetByStatus(context.Background(), status)

		assert.Error(t, err)
	})
}
