package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/billing-engine/internal/commons"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		username := "user123"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`username`,`status`) VALUES (?,?)")).
			WithArgs(username, commons.StatusUserNew).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		user, err := repo.Create(context.Background(), username)
		require.NoError(t, err)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, commons.StatusUserNew, user.Status)
	})

	t.Run("error", func(t *testing.T) {
		username := "user123"

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`username`,`status`) VALUES (?, ?)")).
			WithArgs(username, commons.StatusUserNew).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		_, err := repo.Create(context.Background(), username)
		assert.Error(t, err)
	})
}

func TestUserRepository_GetUser(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		username := "user123"
		row := sqlmock.NewRows([]string{"username", "status"}).
			AddRow(username, commons.StatusUserNew)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ?")).
			WithArgs(username).
			WillReturnRows(row)

		user, err := repo.GetUser(context.Background(), username)
		require.NoError(t, err)
		assert.Equal(t, username, user.Username)
		assert.Equal(t, commons.StatusUserNew, user.Status)
	})

	t.Run("error", func(t *testing.T) {
		username := "user123"

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ?")).
			WithArgs(username).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := repo.GetUser(context.Background(), username)
		assert.Error(t, err)
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		username := "user123"
		status := 2

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `status`=? WHERE username = ?")).
			WithArgs(status, username).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.UpdateUser(context.Background(), username, status)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		username := "user123"
		status := 2

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `status`=? WHERE username = ?")).
			WithArgs(status, username).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.UpdateUser(context.Background(), username, status)
		assert.Error(t, err)
	})
}
