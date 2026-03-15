package auth

import (
	"fmt"

	"github.com/UGRORF/price-tracker/internal/domain"
	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo   *postgres.UserRepo
	jwtService *JWTService
}

func NewService(repo *postgres.UserRepo, service *JWTService) *Service {
	return &Service{
		userRepo:   repo,
		jwtService: service,
	}
}

func (s *Service) Register(username, password string) (*domain.User, string, error) {
	if len(password) < 8 {
		return nil, "", fmt.Errorf("Password should be more 8 symbols")
	}

	if user, err := s.userRepo.GetByUsername(username); user != nil && err == nil {
		return nil, "", fmt.Errorf("User: %s already exists", user.Username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("Error with hashed password: %w", err)
	}

	newUser := &domain.User{
		Username: username,
		Password: string(hashedPassword),
		Role:     domain.RoleUser,
	}

	err = s.userRepo.Create(newUser)
	if err != nil {
		return nil, "", err
	}

	token, err := s.jwtService.CreateToken(newUser)
	if err != nil {
		return nil, "", err
	}

	return newUser, token, nil
}

func (s *Service) Login(username, password string) (*domain.User, string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.jwtService.CreateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
