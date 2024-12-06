package entity

import "time"

type User struct {
	Id         int       `db:"id"`
	Username   string    `db:"username"`
	Password   string    `db:"password"`
	Created_at time.Time `db:"created_at"`
}
