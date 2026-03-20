package domain

import "context"

type User struct {
	ID        string `bun:"id,pk,type:uuid,default:gen_random_uuid()" db:"id"`
	Email     string `bun:"email,notnull,unique" db:"email"`
	Password  string `bun:"password,notnull" db:"password"`
	Role      string `bun:"role,notnull" db:"role"`
	CreatedAt string `bun:"created_at,nullzero,notnull,default:current_timestamp" db:"created_at"`
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, user *User) (User, error)
}
