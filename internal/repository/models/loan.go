package models

type LoanModel struct {
	Id        int     `db:"id"`
	Username  string  `db:"username"`
	Amount    float64 `db:"amount"`
	CreatedAt string  `db:"created_at"`
	Status    int     `db:"status"`
}
