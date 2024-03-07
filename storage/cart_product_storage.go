package storage

import (
	"context"
	"food-menu/types"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type getAllProducts struct {
	CartProducts  []*types.CartProduct `json:"cartProducts"`
	TotalPrice    int                  `json:"totalPrice"`
	TotalQuantity int                  `json:"totalQuantity"`
	Err           error                `json:"-"`
}

type CartProductStorage interface {
	GetAll(userId int) *getAllProducts
	//Create(userId, productId int) error
	//Update(userId, cartProductId, quantity int) error
	//Delete(userId, cartProductId int) error
}

type CartProductPgStorage struct {
	DB *pgxpool.Pool
}

func (c *CartProductPgStorage) GetAll(userId int) *getAllProducts {
	productsCh := make(chan struct {
		products []*types.CartProduct
		err      error
	})
	productsInfoCh := make(chan struct {
		totalPrice    int
		totalQuantity int
		err           error
	})

	go func() {
		data := []*types.CartProduct{}

		if err := pgxscan.Select(context.Background(), c.DB, &data, `WITH p AS (SELECT id as product_id,title,image,price FROM products)
		 SELECT id,p.title,p.price,p.image,quantity,p.product_id FROM cart_product c
				JOIN p ON c.product_id = p.product_id
		WHERE cart_id = (SELECT id FROM cart WHERE user_id = $1);
`, userId); err != nil {
			productsCh <- struct {
				products []*types.CartProduct
				err      error
			}{products: data, err: err}
		} else {
			productsCh <- struct {
				products []*types.CartProduct
				err      error
			}{products: data, err: nil}
		}
	}()
	go func() {

		var info struct {
			TotalPrice int `db:"total_price"`
			Quantity   int `db:"quantity"`
		}

		if err := pgxscan.Get(context.Background(), c.DB, &info, `WITH p AS (SELECT id, price FROM products)
			 SELECT sum(p.price * quantity) as total_price, sum(quantity) as quantity FROM cart_product c
						JOIN p ON c.product_id = p.id
			 WHERE cart_id = (SELECT id FROM cart WHERE user_id = $1);
	`, userId); err != nil {
			productsInfoCh <- struct {
				totalPrice    int
				totalQuantity int
				err           error
			}{totalPrice: 0, totalQuantity: 0, err: err}
		} else {
			productsInfoCh <- struct {
				totalPrice    int
				totalQuantity int
				err           error
			}{totalPrice: info.TotalPrice, totalQuantity: info.Quantity, err: nil}
		}

	}()

	products, productsInfo := <-productsCh, <-productsInfoCh

	data := &getAllProducts{CartProducts: products.products, TotalPrice: productsInfo.totalPrice, TotalQuantity: productsInfo.totalQuantity}
	if products.err != nil {
		data.Err = products.err
	} else {
		data.Err = productsInfo.err
	}

	return data
}
