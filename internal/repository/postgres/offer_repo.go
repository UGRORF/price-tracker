package postgres

import (
	"database/sql"
	"fmt"

	"github.com/UGRORF/price-tracker/internal/domain"
)

type OfferRepo struct {
	db *sql.DB
}

func NewOfferRepo(db *sql.DB) *OfferRepo {
	return &OfferRepo{db: db}
}

// TODO: доделать вставку цены напрямую из ссылки на магазин(карточку товара)
func (r *OfferRepo) Create(offer *domain.Offer) error {
	query := `
		INSERT INTO offers (product_id, store_id, price)
		VALUES ($1, $2, $3)
		RETURNING id
		`

	err := r.db.QueryRow(query, offer.ProductID, offer.StoreID, offer.Price).Scan(&offer.ID)
	if err != nil {
		return fmt.Errorf("Error with create offer: %w", err)
	}

	return nil
}

func (r *OfferRepo) GetByID(id int64) (*domain.Offer, error) {
	query := `
		SELECT 
            o.id, o.product_id, o.store_id, o.price,
            p.id, p.name, p.description,
			s.id, s.name, s.url
        FROM offers o
        JOIN products p ON o.product_id = p.id
        JOIN stores s ON o.store_id = s.id
        WHERE o.id = $1
		`

	offer := &domain.Offer{}
	product := &domain.Product{}
	store := &domain.Store{}

	err := r.db.QueryRow(query, id).Scan(
		&offer.ID, &offer.ProductID, &offer.StoreID, &offer.Price,
		&product.ID, &product.Name, &product.Description,
		&store.ID, &store.Name, &store.URL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("Failed to get offer from DataBase: %w", err)
	}

	offer.Product = product
	offer.Store = store

	return offer, nil
}

func (r *OfferRepo) GetAll() ([]domain.Offer, error) {
	query := `
		SELECT 
            o.id, o.product_id, o.store_id, o.price,
            p.id, p.name, p.description,
			s.id, s.name, s.url
        FROM offers o
        JOIN products p ON o.product_id = p.id
        JOIN stores s ON o.store_id = s.id
		`

	var offers []domain.Offer
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get offers from DataBase: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		offer := domain.Offer{}
		product := &domain.Product{}
		store := &domain.Store{}
		err := rows.Scan(&offer.ID, &offer.ProductID, &offer.StoreID, &offer.Price,
			&product.ID, &product.Name, &product.Description,
			&store.ID, &store.Name, &store.URL)
		if err != nil {
			return nil, fmt.Errorf("Error iterating offers %w", err)
		}

		offer.Product = product
		offer.Store = store
		offers = append(offers, offer)
	}

	return offers, nil
}
