package handlers

import (
	"encoding/json"
	"fmt"
	"food-menu/storage"
	"food-menu/types"
	"net/http"
)

type CartProductHandler struct {
	CartProductStorage storage.CartProductStorage
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
