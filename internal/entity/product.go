package entity

import "time"

type Product struct {
	Id        int        `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"` // Nullable field
	DeletedAt *time.Time `db:"deleted_at"` // Nullable field
}
