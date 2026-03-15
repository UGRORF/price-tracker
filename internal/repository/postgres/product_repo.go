package postgres

import (
	"database/sql"
	"fmt"

	"github.com/UGRORF/price-tracker/internal/domain"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(product *domain.Product) error {
	query := `
		INSERT INTO products (name, description)
		VALUES ($1, $2) 
		RETURNING id
		`

	err := r.db.QueryRow(query, product.Name, product.Description).Scan(&product.ID)
	if err != nil {
		return fmt.Errorf("Failed to create product in DataBase: %w", err)
	}

	return nil
}

func (r *ProductRepo) GetById(id int64) (*domain.Product, error) {
	query := `
		SELECT id, name, description
		FROM products
		WHERE id = $1
		`
	product := &domain.Product{}
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Description)
	if err != nil {
		return nil, fmt.Errorf("Failed to get product from DataBase: %w", err)
	}

	return product, nil
}

func (r *ProductRepo) GetByName(name string) (*domain.Product, error) {
	query := `
		SELECT id, name, description
		FROM products
		WHERE name = $1
		`
	product := &domain.Product{}
	err := r.db.QueryRow(query, name).Scan(&product.ID, &product.Name, &product.Description)
	if err != nil {
		return nil, fmt.Errorf("Failed to get product from DataBase: %w", err)
	}

	return product, nil
}

func (r *ProductRepo) GetAll() ([]domain.Product, error) {
	query := `
		SELECT id, name, description
		FROM products
		`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get products from DataBase: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description)
		if err != nil {
			return nil, fmt.Errorf("Error iterating products %w", err)
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}
