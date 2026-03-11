package domain

type Bucket struct {
	ID          int64    `db:"id" json:"id"`
	UserID      int64    `db:"user_id" json:"user_id"`
	OfferID     int64    `db:"offer_id" json:"offer_id"`
	TargetPrice *float64 `db:"target_price" json:"targetPrice"`

	Offer *Offer `db:"-" json:"offer,omitempty"`
	User  *User  `db:"-" json:"user,omitempty"`
}
