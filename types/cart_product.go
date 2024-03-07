package types

type CartProduct struct {
	ID        int    `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Price     int    `json:"price" db:"price"`
	Image     string `json:"image" db:"image"`
	ProductID int    `json:"productId" db:"product_id"`
	Quantity  int    `json:"quantity" db:"quantity"`
}
