package postgres

import (
	"database/sql"
	"fmt"

	"github.com/UGRORF/price-tracker/internal/domain"
)

type StoreRepo struct {
	db *sql.DB
}

func NewStoreRepo(db *sql.DB) *StoreRepo {
	return &StoreRepo{db: db}
}

func (r *StoreRepo) Create(store *domain.Store) error {
	query := `
		INSERT INTO stores (name, url)
		VALUES ($1, $2) 
		RETURNING id
		`

	err := r.db.QueryRow(query, store.Name, store.URL).Scan(&store.ID)
	if err != nil {
		return fmt.Errorf("Failed to create store in DataBase: %w", err)
	}

	return nil
}

func (r *StoreRepo) GetById(id int64) (*domain.Store, error) {
	query := `
		SELECT id, name, url
		FROM stores
		WHERE id = $1
		`
	store := &domain.Store{}
	err := r.db.QueryRow(query, id).Scan(&store.ID, &store.Name, &store.URL)
	if err != nil {
		return nil, fmt.Errorf("Failed to get store from DataBase: %w", err)
	}

	return store, nil
}

func (r *StoreRepo) GetByName(name string) (*domain.Store, error) {
	query := `
		SELECT id, name, url
		FROM stores
		WHERE name = $1
		`
	store := &domain.Store{}
	err := r.db.QueryRow(query, name).Scan(&store.ID, &store.Name, &store.URL)
	if err != nil {
		return nil, fmt.Errorf("Failed to get store from DataBase: %w", err)
	}

	return store, nil
}

func (r *StoreRepo) GetAll() ([]domain.Store, error) {
	query := `
		SELECT id, name, url
		FROM stores
		`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get stores from DataBase: %w", err)
	}
	defer rows.Close()

	var stores []domain.Store
	for rows.Next() {
		var store domain.Store
		err := rows.Scan(&store.ID, &store.Name, &store.URL)
		if err != nil {
			return nil, fmt.Errorf("Error iterating stores %w", err)
		}
		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stores: %w", err)
	}

	return stores, nil
}
