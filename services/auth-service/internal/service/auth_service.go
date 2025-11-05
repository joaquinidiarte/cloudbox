package service

import (
	"context"
	"errors"

	"github.com/joaquinidiarte/cloudbox/services/auth-service/internal/repository"
	"github.com/joaquinidiarte/cloudbox/shared/models"
	"github.com/joaquinidiarte/cloudbox/shared/utils"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.LoginResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Check if username already exists
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.NewUser(models.UserCreateRequest{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}, hashedPassword)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, expiresAt, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, expiresAt, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (*utils.Claims, error) {
	claims, err := s.jwtManager.Verify(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}
	return claims, nil
}
