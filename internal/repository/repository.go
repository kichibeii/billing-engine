package repository

type Repository struct {
	Loan    ILoanRepository
	User    IUserRepository
	PayLoan IPayLoanRepository
}
