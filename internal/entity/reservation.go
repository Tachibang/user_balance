package entity

import "time"

type Reservation struct {
	Id        int        `db:"id"`
	AccountId int        `db:"account_id"`
	ProductId int        `db:"product_id"`
	Amount    int        `db:"amount"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"` // Nullable field
	DeletedAt *time.Time `db:"deleted_at"` // Nullable field
}
