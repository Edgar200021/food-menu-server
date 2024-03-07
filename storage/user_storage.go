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
	Create(createUser types.CreateUser) (int, error)
}

type UserPgStorage struct {
	DB *pgxpool.Pool
}

func (u *UserPgStorage) Create(createUser types.CreateUser) (int, error) {

	user, err := u.GetByEmail(createUser.Email)
	if err != nil {
		return 0, errors.New(err.Error())
	}

	if user.Email != "" {
		return 0, fmt.Errorf("user with email %s already exists", createUser.Email)
	}

	var userId int
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(createUser.Password), 10)
	if err := u.DB.QueryRow(context.Background(), "INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id", createUser.Email, hashedPassword, createUser.Name).Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
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
