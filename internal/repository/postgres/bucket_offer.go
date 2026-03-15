package postgres

import (
	"database/sql"
	"fmt"

	"github.com/UGRORF/price-tracker/internal/domain"
)

type BucketRepo struct {
	db *sql.DB
}

func NewBucketRepo(db *sql.DB) *BucketRepo {
	return &BucketRepo{db: db}
}

func (r *BucketRepo) GetByUserID(bucketID, userID int64) (*domain.Bucket, error) {
	query := `
		SELECT 
    		b.id, b.user_id, b.offer_id, b.target_price,
			u.id, u.username, u.password,
			o.id, o.product_id, o.store_id, o.price,
			p.id, p.name, p.description,
			s.id, s.name, s.url
		FROM product_buckets b
		JOIN users u ON b.user_id = u.id
		JOIN offers o ON b.offer_id = o.id
		JOIN products p ON o.product_id = p.id
		JOIN stores s ON o.store_id = s.id
		WHERE b.id = $1 and b.user_id = $2
		`

	bucket := &domain.Bucket{}
	offer := &domain.Offer{}
	user := &domain.User{}
	product := &domain.Product{}
	store := &domain.Store{}

	err := r.db.QueryRow(query, bucketID, userID).Scan(
		&bucket.ID, &bucket.UserID, &bucket.OfferID, &bucket.TargetPrice,
		&user.ID, &user.Username, &user.Password,
		&offer.ID, &offer.ProductID, &offer.StoreID, &offer.Price,
		&product.ID, &product.Name, &product.Description,
		&store.ID, &store.Name, &store.URL)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("Failed to get bucket from DataBase: %w", err)
	}

	offer.Product = product
	offer.Store = store
	bucket.Offer = offer
	bucket.User = user

	return bucket, nil
}

func (r *BucketRepo) GetByID(id int64) (*domain.Bucket, error) {
	query := `
		SELECT 
    		b.id, b.user_id, b.offer_id, b.target_price,
			u.id, u.username, u.password,
			o.id, o.product_id, o.store_id, o.price,
			p.id, p.name, p.description,
			s.id, s.name, s.url
		FROM product_buckets b
		JOIN users u ON b.user_id = u.id
		JOIN offers o ON b.offer_id = o.id
		JOIN products p ON o.product_id = p.id
		JOIN stores s ON o.store_id = s.id
		WHERE b.id = $1
		`

	bucket := &domain.Bucket{}
	offer := &domain.Offer{}
	user := &domain.User{}
	product := &domain.Product{}
	store := &domain.Store{}

	err := r.db.QueryRow(query, id).Scan(
		&bucket.ID, &bucket.UserID, &bucket.OfferID, &bucket.TargetPrice,
		&user.ID, &user.Username, &user.Password,
		&offer.ID, &offer.ProductID, &offer.StoreID, &offer.Price,
		&product.ID, &product.Name, &product.Description,
		&store.ID, &store.Name, &store.URL)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("Failed to get bucket from DataBase: %w", err)
	}

	offer.Product = product
	offer.Store = store
	bucket.Offer = offer
	bucket.User = user

	return bucket, nil
}

func (r *BucketRepo) GetAllByUserID(id int64) ([]domain.Bucket, error) {
	query := `
		SELECT 
    		b.id, b.user_id, b.offer_id, b.target_price,
			u.id, u.username, u.password,
			o.id, o.product_id, o.store_id, o.price,
			p.id, p.name, p.description,
			s.id, s.name, s.url
		FROM product_buckets b
		JOIN users u ON b.user_id = u.id
		JOIN offers o ON b.offer_id = o.id
		JOIN products p ON o.product_id = p.id
		JOIN stores s ON o.store_id = s.id
		WHERE user_id = $1
		`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get buckets from DataBase: %w", err)
	}
	defer rows.Close()

	var buckets []domain.Bucket

	for rows.Next() {
		bucket := domain.Bucket{}
		offer := &domain.Offer{}
		user := &domain.User{}
		product := &domain.Product{}
		store := &domain.Store{}

		err := rows.Scan(
			&bucket.ID, &bucket.UserID, &bucket.OfferID, &bucket.TargetPrice,
			&user.ID, &user.Username, &user.Password,
			&offer.ID, &offer.ProductID, &offer.StoreID, &offer.Price,
			&product.ID, &product.Name, &product.Description,
			&store.ID, &store.Name, &store.URL)

		if err != nil {
			return nil, fmt.Errorf("Failed to get bucket from DataBase: %w", err)
		}
		offer.Product = product
		offer.Store = store
		bucket.Offer = offer
		bucket.User = user

		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (r *BucketRepo) GetAll() ([]domain.Bucket, error) {
	query := `
		SELECT 
    		b.id, b.user_id, b.offer_id, b.target_price,
			u.id, u.username, u.password,
			o.id, o.product_id, o.store_id, o.price,
			p.id, p.name, p.description,
			s.id, s.name, s.url
		FROM product_buckets b
		JOIN users u ON b.user_id = u.id
		JOIN offers o ON b.offer_id = o.id
		JOIN products p ON o.product_id = p.id
		JOIN stores s ON o.store_id = s.id
		`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get buckets from DataBase: %w", err)
	}
	defer rows.Close()

	var buckets []domain.Bucket

	for rows.Next() {
		bucket := domain.Bucket{}
		offer := &domain.Offer{}
		user := &domain.User{}
		product := &domain.Product{}
		store := &domain.Store{}

		err := rows.Scan(
			&bucket.ID, &bucket.UserID, &bucket.OfferID, &bucket.TargetPrice,
			&user.ID, &user.Username, &user.Password,
			&offer.ID, &offer.ProductID, &offer.StoreID, &offer.Price,
			&product.ID, &product.Name, &product.Description,
			&store.ID, &store.Name, &store.URL)

		if err != nil {
			return nil, fmt.Errorf("Failed to get bucket from DataBase: %w", err)
		}
		offer.Product = product
		offer.Store = store
		bucket.Offer = offer
		bucket.User = user

		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (r *BucketRepo) Create(bucket *domain.Bucket) error {
	query := `
        INSERT INTO product_buckets (user_id, offer_id, target_price)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err := r.db.QueryRow(query, bucket.UserID, bucket.OfferID, bucket.TargetPrice).Scan(&bucket.ID)
	if err != nil {
		return fmt.Errorf("failed to create bucket in database: %w", err)
	}

	return nil
}
