package handlers

import (
	"encoding/json"
	"fmt"
	"food-menu/storage"
	"food-menu/types"
	"food-menu/utils"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductStorage storage.ProductPgStorage
}

func (p *ProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Provide valid id", http.StatusBadRequest)
		return
	}

	if product, err := p.ProductStorage.Get(int(id)); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	} else {
		if product.ID == 0 {
			w.Write([]byte(fmt.Sprintf("Product with id %d doesn't exists", id)))
			return
		}

		jsonData, err := json.Marshal(product)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}
func (p *ProductHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	if products, err := p.ProductStorage.GetAll(title); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	} else {

		jsonData, err := json.Marshal(products)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}
func (p *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fileName, err := utils.StoreMultipartImage(r, 1024, "image")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	dict := make(map[string]any, 4)
	dict["image"] = "uploads/" + fileName

	for k, val := range r.Form {
		if k == "price" {
			val, _ := strconv.ParseInt(val[0], 10, 32)
			dict[k] = int(val)
		} else if k == "ingredients" {
			dict[k] = val
		} else {
			dict[k] = val[0]

		}
	}

	jsonData, jsonErr := json.Marshal(dict)
	if jsonErr != nil {
		fmt.Println(jsonErr)

		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var createProduct types.CreateProduct

	if err := json.Unmarshal(jsonData, &createProduct); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if validationResult := createProduct.Validate(); len(validationResult) != 0 {
		fmt.Println(validationResult)
		encoded, _ := json.Marshal(validationResult)
		w.Write(encoded)
		return
	}

	if err := p.ProductStorage.Create(&createProduct); err != nil {
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
