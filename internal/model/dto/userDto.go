package dto

import (
	"fmt"
	"go-boilerplate/internal/entity"
	"go-boilerplate/internal/model"
	"time"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Avatar: model.AvatarResponse{
			OriginalName: user.Avatar.OriginalName,
			Key:          user.Avatar.Key,
			Url:          avatarUrl(user.Avatar.Key),
			Size:         user.Avatar.Size,
			MimeType:     user.Avatar.MimeType,
			Visibility:   user.Avatar.Visibility,
		},
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

func avatarUrl(key string) string {
	if key == "" {
		return ""
	}
	return fmt.Sprintf("http://localhost:3000/%s", key)
}
