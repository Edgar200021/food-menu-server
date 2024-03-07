package handlers

import (
	"encoding/json"
	"fmt"
	"food-menu/storage"
	"food-menu/types"
	"io"
	"math"
	"net/http"
	"strconv"
)

type CartProductHandler struct {
	CartProductStorage storage.CartProductStorage
	ProductStorage     storage.ProductStorage
}

func (c *CartProductHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(types.User)

	if data := c.CartProductStorage.GetAll(user.ID); data.Err != nil {
		fmt.Println(data.Err)

		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	} else {
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}
func (c *CartProductHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(types.User)
	productId, convErr := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if convErr != nil {
		http.Error(w, "Provide valid id", http.StatusBadRequest)
		return
	}

	product, err := c.ProductStorage.Get(int(productId))
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if product.ID == 0 {
		http.Error(w, fmt.Sprintf("No product with id %d", productId), http.StatusBadRequest)
		return
	}

	if err := c.CartProductStorage.Create(user.ID, int(productId)); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (c *CartProductHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(types.User)
	productId, convErr := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if convErr != nil {
		http.Error(w, "Provide valid id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decodedBody := make(map[string]int, 1)
	if err := json.Unmarshal(body, &decodedBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if decodedBody["quantity"] == 0 {
		http.Error(w, "Provide quantity of products", http.StatusBadRequest)
		return
	}

	if decodedBody["quantity"] < 0 {
		if cartProduct, err := c.CartProductStorage.Get(user.ID, int(productId)); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		} else {
			if cartProduct.Quantity <= int(math.Abs(float64(decodedBody["quantity"]))) {
				http.Error(w, "Invalid quantity", http.StatusBadRequest)
				return
			}
		}
	}

	if updatedQuantity, err := c.CartProductStorage.Update(user.ID, int(productId), decodedBody["quantity"]); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	} else {
		data := struct {
			Quantity int `json:"quantity"`
		}{Quantity: updatedQuantity}

		encodedData, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Write(encodedData)
	}
}
func (c *CartProductHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(types.User)
	productId, convErr := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if convErr != nil {
		http.Error(w, "Provide valid id", http.StatusBadRequest)
		return
	}

	if err := c.CartProductStorage.Delete(user.ID, int(productId)); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

}
