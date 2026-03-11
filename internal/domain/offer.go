package domain

type Offer struct {
	ID        int64   `db:"id" json:"id"`
	ProductID int64   `db:"product_id" json:"productID"`
	StoreID   int64   `db:"store_id" json:"storeID"`
	Price     float64 `db:"price" json:"price"`

	Product *Product `db:"-" json:"product,omitempty"`
	Store   *Store   `db:"-" json:"store,omitempty"`
}
