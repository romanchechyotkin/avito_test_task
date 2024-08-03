package entity

import "time"

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	UserType  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
