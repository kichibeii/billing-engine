package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/billing-engine/internal/commons"
	mock_repositories "github.com/billing-engine/internal/mock/repository"
	"github.com/billing-engine/internal/repository"
	"github.com/billing-engine/internal/repository/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_IsDelinquent(t *testing.T) {
	t.Run("success get status delinquent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		username := "user123"

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)

		service := NewService(
			&repository.Repository{
				User: userRepoMock,
			},
		)

		userRepoMock.EXPECT().GetUser(gomock.Any(), username).Return(entity.UserEntity{
			Username: username,
			Status:   commons.StatusUserDeliquent,
		}, nil)

		isDelinquent, err := service.IsDelinquent(context.Background(), username)

		assert.Nil(t, err)
		assert.NotNil(t, isDelinquent)
		assert.Equal(t, isDelinquent, true)
	})

	t.Run("success get status not delinquent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		username := "user123"

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)

		service := NewService(
			&repository.Repository{
				User: userRepoMock,
			},
		)

		userRepoMock.EXPECT().GetUser(gomock.Any(), username).Return(entity.UserEntity{
			Username: username,
			Status:   commons.StatusUserNew,
		}, nil)

		isDelinquent, err := service.IsDelinquent(context.Background(), username)

		assert.Nil(t, err)
		assert.NotNil(t, isDelinquent)
		assert.Equal(t, isDelinquent, false)
	})

	t.Run("got error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		username := "user123"

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)

		service := NewService(
			&repository.Repository{
				User: userRepoMock,
			},
		)

		userRepoMock.EXPECT().GetUser(gomock.Any(), username).Return(entity.UserEntity{}, errors.New("any error from repository"))

		_, err := service.IsDelinquent(context.Background(), username)

		assert.NotNil(t, err)

	})
}

func TestService_CreateLoan(t *testing.T) {
	t.Run("success craete loan not new user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := CreateLoanEntity{
			Username: "user123",
			Amount:   50000000,
		}

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		userRepoMock.EXPECT().GetUser(gomock.Any(), data.Username).Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserNew,
		}, nil)

		loaRepoMock.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(entity.LoanEntity{
			Id:        1,
			Username:  "user123",
			Amount:    50000000,
			Status:    0,
			CreatedAt: time.Now(),
		}, nil)

		payLoanRepoMock.EXPECT().BatchInsert(gomock.Any(), gomock.Any()).Return(nil)

		userRepoMock.EXPECT().UpdateUser(gomock.Any(), "user123", commons.StatusUserActiveLoan).Return(nil)

		err := service.CreateLoan(context.Background(), data)

		assert.Nil(t, err)
	})

	t.Run("success craete loan new user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := CreateLoanEntity{
			Username: "user123",
			Amount:   50000000,
		}

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		userRepoMock.EXPECT().GetUser(gomock.Any(), data.Username).Return(entity.UserEntity{}, nil)

		userRepoMock.EXPECT().Create(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserNew,
		}, nil)

		loaRepoMock.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(entity.LoanEntity{
			Id:        1,
			Username:  "user123",
			Amount:    50000000,
			Status:    0,
			CreatedAt: time.Now(),
		}, nil)

		payLoanRepoMock.EXPECT().BatchInsert(gomock.Any(), gomock.Any()).Return(nil)

		userRepoMock.EXPECT().UpdateUser(gomock.Any(), "user123", commons.StatusUserActiveLoan).Return(nil)

		err := service.CreateLoan(context.Background(), data)

		assert.Nil(t, err)
	})

	t.Run("error user already have loan", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := CreateLoanEntity{
			Username: "user123",
			Amount:   50000000,
		}

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		userRepoMock.EXPECT().GetUser(gomock.Any(), data.Username).Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserActiveLoan,
		}, nil)

		err := service.CreateLoan(context.Background(), data)

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "user have other active loan")
	})
}

func TestService_GetOutstanding(t *testing.T) {
	t.Run("succes get outstanding", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		userRepoMock.EXPECT().GetUser(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserActiveLoan,
		}, nil)

		loaRepoMock.EXPECT().Get(gomock.Any(), "user123", commons.StatusLoanNew).Return(entity.LoanEntity{
			Id:        123,
			Username:  "user123",
			Amount:    55000000,
			Status:    commons.StatusLoanNew,
			CreatedAt: time.Now(),
		}, nil)

		payLoanRepoMock.EXPECT().GetPayLoanByLoanId(gomock.Any(), gomock.Any()).Return([]entity.PayLoanEntity{
			{
				Id:     123,
				Amount: 550000,
				Status: commons.StatusPayLoanPayed,
			}, {
				Id:     124,
				Amount: 550000,
				Status: commons.StatusPayLoanUnpayed,
			},
		}, nil)

		amount, err := service.GetOutStanding(context.Background(), "user123")

		assert.Nil(t, err)
		assert.NotNil(t, amount)
		assert.Equal(t, amount, float64(54450000))
	})

	t.Run("error when user not active loan status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		userRepoMock.EXPECT().GetUser(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserClosedLoan,
		}, nil)

		amount, err := service.GetOutStanding(context.Background(), "user123")

		assert.NotNil(t, err)
		assert.Equal(t, amount, float64(0))
	})
}

func TestService_MakePayment(t *testing.T) {
	t.Run("success make payment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		data := MakePaymentEntity{
			Username: "user123",
			Amount:   5500000,
		}

		userRepoMock.EXPECT().GetUser(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserActiveLoan,
		}, nil)

		loaRepoMock.EXPECT().Get(gomock.Any(), "user123", commons.StatusLoanNew).Return(entity.LoanEntity{
			Id:        123,
			Username:  "user123",
			Amount:    55000000,
			Status:    commons.StatusLoanNew,
			CreatedAt: time.Now(),
		}, nil)

		payLoanRepoMock.EXPECT().GetInSpecificTimeAndStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return([]entity.PayLoanEntity{
			{
				Id:     123,
				LoanId: 123,
				Amount: 5500000,
				Status: commons.StatusPayLoanUnpayed,
			},
		}, nil)

		payLoanRepoMock.EXPECT().Update(gomock.Any(), gomock.Any(), entity.PayLoanEntity{
			Status: commons.StatusPayLoanPayed,
		})

		message, err := service.MakePayment(context.Background(), data)

		assert.Nil(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, message, "success make payment")
	})

	t.Run("error amount not same with requirment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		data := MakePaymentEntity{
			Username: "user123",
			Amount:   500000,
		}

		userRepoMock.EXPECT().GetUser(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserActiveLoan,
		}, nil)

		loaRepoMock.EXPECT().Get(gomock.Any(), "user123", commons.StatusLoanNew).Return(entity.LoanEntity{
			Id:        123,
			Username:  "user123",
			Amount:    55000000,
			Status:    commons.StatusLoanNew,
			CreatedAt: time.Now(),
		}, nil)

		payLoanRepoMock.EXPECT().GetInSpecificTimeAndStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return([]entity.PayLoanEntity{
			{
				Id:     123,
				LoanId: 123,
				Amount: 5500000,
				Status: commons.StatusPayLoanUnpayed,
			},
		}, nil)

		message, err := service.MakePayment(context.Background(), data)

		assert.Nil(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, message, fmt.Sprintf("amount not same with requirment : %.2f", float64(5500000)))
	})

	t.Run("error already pay", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		data := MakePaymentEntity{
			Username: "user123",
			Amount:   500000,
		}

		userRepoMock.EXPECT().GetUser(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserActiveLoan,
		}, nil)

		loaRepoMock.EXPECT().Get(gomock.Any(), "user123", commons.StatusLoanNew).Return(entity.LoanEntity{
			Id: 123,
		}, nil)

		payLoanRepoMock.EXPECT().GetInSpecificTimeAndStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return([]entity.PayLoanEntity{}, nil)

		message, err := service.MakePayment(context.Background(), data)

		assert.Nil(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, message, "already payed for this week")
	})

	t.Run("error not in active loan", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		data := MakePaymentEntity{
			Username: "user123",
			Amount:   500000,
		}

		userRepoMock.EXPECT().GetUser(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserActiveLoan,
		}, nil)

		loaRepoMock.EXPECT().Get(gomock.Any(), "user123", commons.StatusLoanNew).Return(entity.LoanEntity{}, nil)

		message, err := service.MakePayment(context.Background(), data)

		assert.Nil(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, message, "not have any active loan")
	})

	t.Run("user not have active loan", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		data := MakePaymentEntity{
			Username: "user123",
			Amount:   500000,
		}

		userRepoMock.EXPECT().GetUser(gomock.Any(), "user123").Return(entity.UserEntity{
			Username: "user123",
			Status:   commons.StatusUserClosedLoan,
		}, nil)

		message, err := service.MakePayment(context.Background(), data)

		assert.Nil(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, message, "user not active loan")
	})
}

func TestService_ScheduleTask(t *testing.T) {
	t.Run("success flow schedule", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		userRepoMock := mock_repositories.NewMockIUserRepository(ctrl)
		loaRepoMock := mock_repositories.NewMockILoanRepository(ctrl)
		payLoanRepoMock := mock_repositories.NewMockIPayLoanRepository(ctrl)

		service := NewService(&repository.Repository{
			User:    userRepoMock,
			Loan:    loaRepoMock,
			PayLoan: payLoanRepoMock,
		})

		loaRepoMock.EXPECT().GetByStatus(gomock.Any(), commons.StatusLoanNew).Return([]entity.LoanEntity{
			{
				Id:        123,
				Username:  "bambang1",
				Amount:    55000000,
				Status:    commons.StatusLoanNew,
				CreatedAt: time.Now(),
			},
			{
				Id:        124,
				Username:  "bambang2",
				Amount:    55000000,
				Status:    commons.StatusLoanNew,
				CreatedAt: time.Now(),
			},
		}, nil)

		payLoanRepoMock.EXPECT().GetInSpecificTimeAndStatus(gomock.Any(), 123, gomock.Any()).Return([]entity.PayLoanEntity{}, nil)

		loaRepoMock.EXPECT().UpdateStatus(gomock.Any(), 123, commons.StatusLoanClosed).Return(nil)

		userRepoMock.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), commons.StatusUserClosedLoan).Return(nil)

		payLoanRepoMock.EXPECT().GetInSpecificTimeAndStatus(gomock.Any(), 124, gomock.Any()).Return([]entity.PayLoanEntity{
			{
				Id:     123,
				LoanId: 124,
				Amount: 50000,
				Status: commons.StatusPayLoanUnpayed,
			},
			{
				Id:     124,
				LoanId: 124,
				Amount: 50000,
				Status: commons.StatusPayLoanUnpayed,
			},
			{
				Id:     125,
				LoanId: 124,
				Amount: 50000,
				Status: commons.StatusPayLoanUnpayed,
			},
		}, nil)

		userRepoMock.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), commons.StatusUserDeliquent).Return(nil)

		err := service.ScheduleTask(context.Background())

		assert.Nil(t, err)
	})
}
