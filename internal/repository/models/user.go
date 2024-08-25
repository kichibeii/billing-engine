package models

type UserModel struct {
	Username string `db:"username"`
	Status   int    `db:"status"`
}
