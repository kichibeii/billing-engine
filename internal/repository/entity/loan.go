package entity

import "time"

type LoanEntity struct {
	Id        int
	Username  string
	Amount    float64
	Status    int
	CreatedAt time.Time
}
