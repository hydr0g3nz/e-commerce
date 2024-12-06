package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	accessSecret  string
	refreshSecret string
	repository    ports.AuthRepository
}

func NewAuthService(accessSecret, refreshSecret string, repository ports.AuthRepository) *AuthService {
	return &AuthService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		repository:    repository,
	}
}
func (s *AuthService) GenerateTokenPair(userId string, role string) (*domain.TokenResponse, error) {
	// Generate access token
	accessToken, expiresIn, err := s.createAccessToken(userId, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate refresh token (longer lived)
	refreshToken, err := s.createRefreshToken(userId, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (*domain.TokenResponse, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token claims")
	}

	// Extract token metadata
	userId, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in refresh token")
	}

	tokenUuid, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.New("invalid uuid in refresh token")
	}
	// Verify token exists in storage and hasn't been revoked
	storedToken, err := s.repository.FetchRefreshToken(tokenUuid)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or revoked: %w", err)
	}

	if storedToken.UserID != userId {
		return nil, errors.New("refresh token user mismatch")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role in refresh token")
	}
	// Generate only a new access token
	accessToken, expiresIn, err := s.createAccessToken(userId, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create new access token: %w", err)
	}

	// Return only the new access token
	return &domain.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}, nil
}

func (s *AuthService) createAccessToken(userId string, role string) (string, int64, error) {
	now := time.Now()
	expiresIn := int64(15 * 60) // 15 minutes in seconds
	exp := now.Add(time.Second * time.Duration(expiresIn)).Unix()

	claims := jwt.MapClaims{
		"user_id": userId,
		"uuid":    uuid.New().String(),
		"exp":     exp,
		"iat":     now.Unix(),
		"type":    "access",
		"role":    role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, expiresIn, nil
}

func (s *AuthService) createRefreshToken(userId, role string) (string, error) {
	now := time.Now()
	tokenUuid := uuid.New().String()

	// Refresh token expires in 7 days
	exp := now.Add(time.Hour * 24 * 7).Unix()

	claims := jwt.MapClaims{
		"user_id": userId,
		"uuid":    tokenUuid,
		"exp":     exp,
		"iat":     now.Unix(),
		"role":    role,
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return "", err
	}

	// Store refresh token metadata
	err = s.repository.CreateRefreshToken(userId, &domain.TokenMetadata{
		UserId: userId,
		Uuid:   tokenUuid,
		Exp:    exp,
	})
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return signedToken, nil
}
func (s *AuthService) Login(credentials *domain.UserCredentials) (*domain.User, *domain.TokenResponse, error) {
	// Validate user credentials
	user, err := s.ValidateCredentials(credentials.Email, credentials.Password)
	if err != nil {
		return nil, nil, err
	}

	// Generate token pair
	tokenDetails, err := s.GenerateTokenPair(user.ID, user.Role)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenDetails, nil
}

func (s *AuthService) ValidateCredentials(email, password string) (*domain.User, error) {
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		// Use a generic error message to prevent email enumeration
		return nil, errors.New("invalid credentials")
	}

	if err := s.ValidatePassword(user.Password, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Clear sensitive data before returning
	user.Password = ""
	return user, nil
}
func (r *AuthService) ValidatePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}
func (r *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *AuthService) Register(request *domain.User) (*domain.TokenResponse, error) {
	// Check if email already exists
	if exists := s.repository.EmailExists(request.Email); exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := s.HashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create new user
	user := &domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashedPassword,
	}

	// Save user to database
	if err := s.repository.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	fmt.Printf("user : %+v", user)
	// Generate tokens
	tokenDetails, err := s.GenerateTokenPair(user.ID, "user")
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}
	return tokenDetails, nil
}
