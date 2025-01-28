package dto

import (
	"go-api-fiber/internal/entity"
	"go-api-fiber/internal/model"
	"time"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: time.Time.Format(user.CreatedAt, time.RFC3339),
		UpdatedAt: time.Time.Format(user.CreatedAt, time.RFC3339),
	}
}

func LoginUserToReponse(user *entity.User, token string) *model.LoginUserResponse {
	return &model.LoginUserResponse{
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Token:     token,
		CreatedAt: time.Time.Format(user.CreatedAt, time.RFC3339),
		UpdatedAt: time.Time.Format(user.CreatedAt, time.RFC3339),
	}
}
