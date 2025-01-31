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
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func LoginUserToReponse(user *entity.User, token string, expired time.Time) *model.LoginUserResponse {
	return &model.LoginUserResponse{
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		AccessToken: model.DetailToken{
			Token:   token,
			Expired: expired.Format(time.RFC3339),
		},
	}
}
