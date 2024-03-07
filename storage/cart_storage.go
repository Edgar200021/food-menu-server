package storage

import (
	"context"
	"errors"
	"fmt"
	"food-menu/types"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartStorage interface {
	Get(userId int) (types.Cart, error)
	Create(userId int) error
}

type CartPgStorage struct {
	DB *pgxpool.Pool
}

func (c *CartPgStorage) Get(userId int) (types.Cart, error) {
	var cart types.Cart

	if err := pgxscan.Get(context.Background(), c.DB, &cart, "SELECT * FROM cart WHERE user_id = $1", userId); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return cart, err
		}
	}

	return cart, nil
}

func (c *CartPgStorage) Create(userId int) error {

	if _, err := c.DB.Query(context.Background(), "INSERT INTO cart (user_id) VALUES ($1)", userId); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
