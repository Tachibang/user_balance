package entity

import "time"

type Account struct {
	Id        int        `db:"id"`
	Balance   int        `db:"balance"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"` // Nullable field
	DeletedAt *time.Time `db:"deleted_at"` // Nullable field
}
