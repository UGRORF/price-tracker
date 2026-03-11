package domain

type Store struct {
	ID   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	URL  string `db:"url" json:"URL"`
}
