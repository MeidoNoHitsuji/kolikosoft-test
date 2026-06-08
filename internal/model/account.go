package model

import "time"

type Account struct {
	ID      int64 `db:"id" json:"id"`
	Balance int64 `db:"balance" json:"balance"`
}

type AccountHistory struct {
	ID         int64     `db:"id" json:"id"`
	AccountID  int64     `db:"account_id" json:"account_id"`
	OldBalance int64     `db:"old_balance" json:"old_balance"`
	NewBalance int64     `db:"new_balance" json:"new_balance"`
	Comment    *string   `db:"comment" json:"comment,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
