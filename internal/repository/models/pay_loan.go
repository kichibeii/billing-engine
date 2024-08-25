package models

type PayLoanModel struct {
	Id        int     `db:"id"`
	LoanId    int     `db:"loan_id"`
	Amount    float64 `db:"amount"`
	CreatedAt string  `db:"created_at"`
	Status    int     `db:"status"`
}
