package commons

import "time"

// different time
const (
	DifferentTime = 30 * time.Second
)

// interest
const (
	FlatInterest = 10
)

// week pay
const (
	CountWeekPay = 50
)

// status user
const (
	StatusUserNew        = 1
	StatusUserActiveLoan = 2
	StatusUserDeliquent  = 3
	StatusUserClosedLoan = 4
)

// status loan
const (
	StatusLoanNew    = 0
	StatusLoanClosed = 1
)

// status payloan
const (
	StatusPayLoanUnpayed = 0
	StatusPayLoanPayed   = 1
)
