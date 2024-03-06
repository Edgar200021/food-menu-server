package storage

import (
	"context"
	"errors"
	"fmt"
	"food-menu/types"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	GetByEmail(email string) (types.User, error)
	GetById(id int) (types.User, error)
	Create(createUser types.CreateUser) error
}

type UserPgStorage struct {
	DB *pgxpool.Pool
}

func (u *UserPgStorage) Create(createUser types.CreateUser) error {

	user, err := u.GetByEmail(createUser.Email)
	if err != nil {
		return errors.New(err.Error())
	}

	if user.Email != "" {
		return fmt.Errorf("user with email %s already exists", createUser.Email)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(createUser.Password), 10)
	if _, err := u.DB.Query(context.Background(), "INSERT INTO users (email, password, name) VALUES ($1, $2, $3)", createUser.Email, hashedPassword, createUser.Name); err != nil {
		return err
	}

	return nil
}
func (u *UserPgStorage) GetByEmail(email string) (types.User, error) {
	var user types.User

	if err := pgxscan.Get(context.Background(), u.DB, &user, `SELECT id, name, email, password,avatar
															  FROM users
															  WHERE email = $1`, email); err != nil {

		if !errors.Is(err, pgx.ErrNoRows) {
			return user, err
		}
	}

	return user, nil
}
func (u *UserPgStorage) GetById(id int) (types.User, error) {
	var user types.User

	if err := pgxscan.Get(context.Background(), u.DB, &user, `SELECT id, name, email, password,avatar
															  FROM users
															  WHERE id = $1`, id); err != nil {

		if !errors.Is(err, pgx.ErrNoRows) {
			return user, err
		}
	}

	return user, nil
}
