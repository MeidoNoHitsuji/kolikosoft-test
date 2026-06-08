package model

type AddHistory struct {
	AccountID  int64  `db:"account_id"`
	OldBalance int64  `db:"old_balance"`
	NewBalance int64  `db:"new_balance"`
	Comment    string `db:"-"`
}
