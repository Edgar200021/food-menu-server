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

type ProductStorage interface {
	Get(id int) (types.Product, error)
	GetAll(title string) ([]*types.Product, error)
	Create(product types.CreateProduct) error
	Delete(id int) error
}

type ProductPgStorage struct {
	DB *pgxpool.Pool
}

func (p *ProductPgStorage) Get(id int) (types.Product, error) {
	var product types.Product

	if err := pgxscan.Get(context.Background(), p.DB, &product, `SELECT * FROM products WHERE id = $1`, id); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return product, err
		}
	}

	return product, nil

}
func (p *ProductPgStorage) GetAll(title string) ([]*types.Product, error) {
	query := `SELECT * FROM products`
	var queryParams []interface{}

	if title != "" {
		query += ` WHERE title LIKE $1`
		queryParams = append(queryParams, "%"+title+"%")
	}

	var products = []*types.Product{}

	if err := pgxscan.Select(context.Background(), p.DB, &products, query, queryParams...); err != nil {
		return products, err
	}
	fmt.Println(products)

	return products, nil
}
func (p *ProductPgStorage) Create(product *types.CreateProduct) error {

	if _, err := p.DB.Query(context.Background(), "INSERT INTO products (title,image, price,ingredients) VALUES ($1, $2, $3, $4)", product.Title, product.Image, product.Price, product.Ingredients); err != nil {
		return err
	}

	return nil
}
func (p *ProductPgStorage) Delete(id int) error {

	product, err := p.Get(id)
	if err != nil || product.ID == 0 {
		return err
	}

	if _, err := p.DB.Query(context.Background(), "DELETE FROM products WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
