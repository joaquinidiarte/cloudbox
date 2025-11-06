package service

import (
	"context"

	"github.com/joaquinidiarte/cloudbox/services/user-service/internal/repository"
	"github.com/joaquinidiarte/cloudbox/shared/models"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	if err := s.userRepo.Update(ctx, id, req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := user.ToResponse()
	return &response, nil
}

func (s *UserService) UpdateStorageUsed(ctx context.Context, userID string, delta int64) error {
	return s.userRepo.UpdateStorageUsed(ctx, userID, delta)
}
