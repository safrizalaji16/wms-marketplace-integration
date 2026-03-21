package repository

import (
	"backend/domain"
	"context"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUser(db *bun.DB) domain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	err := u.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) (domain.User, error) {
	err := u.db.NewInsert().Model(user).Returning("*").Scan(ctx)

	if err != nil {
		return domain.User{}, err
	}

	return *user, nil
}
