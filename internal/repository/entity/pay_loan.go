package entity

import "time"

type PayLoanEntity struct {
	Id        int
	LoanId    int
	Amount    float64
	CreatedAt time.Time
	Status    int
}
