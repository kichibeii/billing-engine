package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/billing-engine/internal/commons"
	"github.com/billing-engine/internal/repository"
	"github.com/billing-engine/internal/repository/entity"
)

type CreateLoanEntity struct {
	Username string
	Amount   float64
}

type MakePaymentEntity struct {
	Username string
	Amount   float64
}

type Service struct {
	repo *repository.Repository
}

type ServiceInterface interface {
	ScheduleTask(ctx context.Context) error
	GetOutStanding(ctx context.Context, username string) (float64, error)
	CreateLoan(ctx context.Context, data CreateLoanEntity) error
	IsDelinquent(ctx context.Context, username string) (bool, error)
	MakePayment(ctx context.Context, data MakePaymentEntity) (string, error)
}

func NewService(repo *repository.Repository) ServiceInterface {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ScheduleTask(ctx context.Context) error {
	// check all open loan from users
	openLoans, err := s.repo.Loan.GetByStatus(ctx, commons.StatusLoanNew)
	if err != nil {
		return err
	}

	// check every loan that open / closed pay_loan
	for _, openLoan := range openLoans {
		payLoans, err := s.repo.PayLoan.GetInSpecificTimeAndStatus(ctx, openLoan.Id, time.Now())
		if err != nil {
			return err
		}

		if len(payLoans) == 0 {
			// update loan status to closed
			err = s.repo.Loan.UpdateStatus(ctx, openLoan.Id, commons.StatusLoanClosed)
			if err != nil {
				return err
			}

			// update user status to closed loan
			err = s.repo.User.UpdateUser(ctx, openLoan.Username, commons.StatusUserClosedLoan)
			if err != nil {
				return err
			}
		} else if len(payLoans) >= 2 {
			// update user to delinquent
			err = s.repo.User.UpdateUser(ctx, openLoan.Username, commons.StatusUserDeliquent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) MakePayment(ctx context.Context, data MakePaymentEntity) (string, error) {
	// check user have loan
	user, err := s.repo.User.GetUser(ctx, data.Username)
	if err != nil {
		return "", err
	}

	if user.Status != commons.StatusUserActiveLoan {
		return "user not active loan", nil
	}

	// get loan
	loan, err := s.repo.Loan.Get(ctx, data.Username, commons.StatusLoanNew)
	if err != nil {
		return "", err
	}

	if loan.Id == 0 {
		return "not have any active loan", nil
	}

	payloans, err := s.repo.PayLoan.GetInSpecificTimeAndStatus(ctx, loan.Id, time.Now())
	if err != nil {
		return "", err
	}

	if len(payloans) == 0 {
		return "already payed for this week", nil
	}

	sort.Slice(payloans, func(i, j int) bool {
		return payloans[i].CreatedAt.Before(payloans[j].CreatedAt)
	})

	if payloans[0].Amount != data.Amount {
		return fmt.Sprintf("amount not same with requirment : %.2f", payloans[0].Amount), nil
	}

	err = s.repo.PayLoan.Update(ctx, payloans[0].Id, entity.PayLoanEntity{
		Status: commons.StatusPayLoanPayed,
	})
	if err != nil {
		return "", err
	}

	return "success make payment", nil
}

func (s *Service) IsDelinquent(ctx context.Context, username string) (bool, error) {
	user, err := s.repo.User.GetUser(ctx, username)
	if err != nil {
		return false, err
	}

	return user.Status == commons.StatusUserDeliquent, nil
}

func (s *Service) CreateLoan(ctx context.Context, data CreateLoanEntity) error {
	var user entity.UserEntity
	var err error

	// check user active loan or not
	// validate one user only can make one loan
	user, err = s.repo.User.GetUser(ctx, data.Username)
	if err != nil {
		return err
	}

	if user.Username == "" {
		userCreate, err := s.repo.User.Create(ctx, data.Username)
		if err != nil {
			return err
		}

		user = userCreate
	}

	if user.Status == commons.StatusUserActiveLoan {
		return errors.New("user have other active loan")
	}

	// create loan data
	loan, err := s.repo.Loan.CreateLoan(ctx, entity.LoanEntity{
		Username: data.Username,
		// amount that saved on loan after add interest fee
		Amount:    data.Amount + data.Amount*commons.FlatInterest/100,
		CreatedAt: time.Now(),
		Status:    commons.StatusLoanNew,
	})
	if err != nil {
		return err
	}

	// create pay_loan data for several weeks payment
	amountPerPay := loan.Amount / commons.CountWeekPay
	timePay := loan.CreatedAt
	amountPerPayAfterInterest := amountPerPay*commons.FlatInterest/100 + amountPerPay
	payLoanEntities := []entity.PayLoanEntity{}
	for i := 0; i < commons.CountWeekPay; i++ {
		payLoanEntities = append(payLoanEntities, entity.PayLoanEntity{
			LoanId:    loan.Id,
			Amount:    amountPerPayAfterInterest,
			CreatedAt: timePay,
			Status:    commons.StatusPayLoanUnpayed,
		})

		timePay = timePay.Add(commons.DifferentTime)
	}

	err = s.repo.PayLoan.BatchInsert(ctx, payLoanEntities)
	if err != nil {
		return err
	}

	// update user status
	err = s.repo.User.UpdateUser(ctx, data.Username, commons.StatusUserActiveLoan)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetOutStanding(ctx context.Context, username string) (float64, error) {
	// get users with status loan
	user, err := s.repo.User.GetUser(ctx, username)
	if err != nil {
		return 0, err
	}

	if user.Status != commons.StatusUserActiveLoan {
		return 0, errors.New("user not on open loan")
	}

	// get loan data
	loan, err := s.repo.Loan.Get(ctx, username, commons.StatusLoanNew)
	if err != nil {
		return 0, err
	}

	// get loan_pay
	payLoans, err := s.repo.PayLoan.GetPayLoanByLoanId(ctx, loan.Id)
	if err != nil {
		return 0, err
	}

	// calculate out standing
	payed := float64(0)
	for _, payLoan := range payLoans {
		if payLoan.Status == commons.StatusPayLoanPayed {
			payed += payLoan.Amount
		}
	}
	outstanding := loan.Amount - payed

	return outstanding, nil
}
