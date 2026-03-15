package postgres

import (
	"database/sql"
	"fmt"

	"github.com/UGRORF/price-tracker/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *domain.User) error {
	if !user.IsValidRole() {
		return fmt.Errorf("Not valid role")
	}

	query := `
		INSERT INTO users (username, password, role)
		VALUES ($1, $2, $3) 
		RETURNING id
		`

	err := r.db.QueryRow(query, user.Username, user.Password, user.Role).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("Failed to create user in DataBase: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByID(id int64) (*domain.User, error) {
	query := `
		SELECT id, username, password, role
		from users 
		WHERE id = $1
		`
	user := &domain.User{}

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("Failed to get user from DataBase: %w", err)
	}
	if !user.IsValidRole() {
		return nil, fmt.Errorf("Not valid role")
	}

	return user, nil
}

func (r *UserRepo) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, password, role
		from users
		WHERE username = $1
		`

	user := &domain.User{}
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("Failed to get user from DataBase: %w", err)
	}
	if !user.IsValidRole() {
		return nil, fmt.Errorf("Not valid role")
	}

	return user, nil
}

func (r *UserRepo) GetAll() ([]domain.User, error) {
	query := `
		SELECT id, username, role
		FROM users
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get users from DataBase: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Username, &user.Role)
		if err != nil {
			return nil, fmt.Errorf("error iterating users: %w", err)
		}

		if !user.IsValidRole() {
			continue
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}
