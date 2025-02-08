package dto

import (
	"go-boilerplate/internal/entity"
	"go-boilerplate/internal/model"
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

func LoginUserToReponse(user *entity.User, token string, expiresAt int64) *model.LoginUserResponse {
	timestamp := int64(expiresAt)
	exp := time.Unix(timestamp, 0)

	return &model.LoginUserResponse{
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		AccessToken: model.DetailToken{
			Token:     token,
			ExpiredAt: exp.Format(time.RFC3339),
		},
	}
}
