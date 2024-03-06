package types

import (
	"strings"
)

type Product struct {
	ID          int      `json:"id" db:"id"`
	Title       string   `json:"title" db:"title"`
	Image       string   `json:"image" db:"image"`
	Price       int      `json:"price" db:"price"`
	Ingredients []string `json:"ingredients" db:"ingredients"`
}

type CreateProduct struct {
	Title       string   `json:"title" db:"title"`
	Image       string   `json:"image" db:"image"`
	Price       int      `json:"price" db:"price"`
	Ingredients []string `json:"ingredients" db:"ingredients"`
}

func (c *CreateProduct) Validate() map[string]string {
	fields := make(map[string]string)
	if strings.TrimSpace(c.Title) == "" {
		fields["title"] = "Provide title"
	}

	if strings.TrimSpace(c.Image) == "" {
		fields["image"] = "Provide image"
	}

	if c.Price <= 0 {
		fields["price"] = "Price of product must be more than 0 "
	}

	if len(c.Ingredients) == 0 {
		fields["ingredients"] = "Product must have a ingredients"
	}

	return fields
}
