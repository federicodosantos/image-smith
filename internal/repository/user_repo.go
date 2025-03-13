package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/federicodosantos/image-smith/internal/model"
	"github.com/federicodosantos/image-smith/internal/repository/query"

	customErr "github.com/federicodosantos/image-smith/pkg/error"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserById(ctx context.Context, id string) (*model.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	exists, err := u.CheckEmailExists(ctx, user.Email)
	if err != nil {
		return err
	}

	if exists {
		return customErr.ErrEmailExist
	}

	result, err := u.db.ExecContext(ctx, query.InsertUserQuery,
		user.ID, user.Name, user.Email, user.Password, user.UpdatedAt, user.UpdatedAt)
	if err != nil {
		return nil
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rows != 1 {
		return customErr.ErrRowsAffected
	}

	return nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := u.db.GetContext(ctx, &user, query.GetUserByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErr.ErrEmailNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := u.db.GetContext(ctx, &user, query.GetUserByIdQuery, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	var count int

	err := u.db.GetContext(ctx, &count, query.CheckEmailExistQuery, email)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
