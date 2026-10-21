package model

import "time"

type User struct {
	Id        int       `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
}
