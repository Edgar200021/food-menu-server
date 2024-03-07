package types

type Cart struct {
	ID     int `json:"id" db:"id"`
	UserID int `json:"userId" db:"user_id"`
}
