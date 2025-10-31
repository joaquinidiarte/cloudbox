package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Email     string    `json:"email" bson:"email"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"-" bson:"password"`
	FirstName string    `json:"first_name" bson:"first_name"`
	LastName  string    `json:"last_name" bson:"last_name"`
	StorageUsed int64   `json:"storage_used" bson:"storage_used"`
	StorageLimit int64  `json:"storage_limit" bson:"storage_limit"`
	IsActive  bool      `json:"is_active" bson:"is_active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UserCreateRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty" binding:"omitempty,email"`
}

type UserResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	StorageUsed  int64     `json:"storage_used"`
	StorageLimit int64     `json:"storage_limit"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUser(req UserCreateRequest, hashedPassword string) *User {
	now := time.Now()
	return &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		Username:     req.Username,
		Password:     hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		StorageUsed:  0,
		StorageLimit: 5 * 1024 * 1024 * 1024,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		StorageUsed:  u.StorageUsed,
		StorageLimit: u.StorageLimit,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
	}
}