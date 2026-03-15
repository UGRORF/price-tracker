package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/UGRORF/price-tracker/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey   []byte
	tokenExpiry time.Duration
}

func NewJWTService() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}

	return &JWTService{
		secretKey:   []byte(secret),
		tokenExpiry: 24 * time.Hour,
	}
}

type Claims struct {
	ID       int64
	Username string
	Role     domain.UserRole
	jwt.RegisteredClaims
}

func (s *JWTService) CreateToken(user *domain.User) (string, error) {
	claims := &Claims{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	singedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("Failed to singn token: %w", err)
	}

	return singedToken, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error with validate token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	newClaims := &Claims{
		ID:       claims.ID,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	singedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("Failed to singned token: %w", err)
	}

	return singedToken, nil
}

func (s *JWTService) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.ID, nil
}
